package concernedlogic

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/concerned"
	"rpc-social/pkg/code"
	"rpc-social/pkg/feed"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedListLogic {
	return &ConcernedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConcernedListLogic) ConcernedList(in *social.ConcernedListRequest) (resp *social.ConcernedListResponse, err error) {
	resp = new(social.ConcernedListResponse)
	resp.Data = new(social.ConcernedListData)
	resp.Data.Items = make([]*social.ConcernedItem, 0)

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	// 优先尝试读取推拉结合 Feed 流
	if l.svcCtx.BizRedis != nil && in.BizId == "feed" {
		inboxKey := fmt.Sprintf("biz#feed#inbox:%d", in.UserId)
		pairs, rerr := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(
			l.ctx, inboxKey, 0, in.Cursor, 0, int(in.PageSize)+1)
		if rerr == nil && len(pairs) > 0 {
			var inboxItems []*feed.FeedItem
			for _, pair := range pairs {
				objId, _ := strconv.ParseInt(pair.Key, 10, 64)
				inboxItems = append(inboxItems, &feed.FeedItem{
					ArticleID:   objId,
					PublishTime: pair.Score,
				})
			}
			merger := feed.NewFeedMerger()
			merged := merger.MergeMultiStreams([][]*feed.FeedItem{inboxItems}, int(in.PageSize))
			for _, item := range merged {
				resp.Data.Items = append(resp.Data.Items, &social.ConcernedItem{
					Id:         item.ArticleID,
					BizId:      "feed",
					ObjId:      item.ArticleID,
					CreateTime: item.PublishTime,
				})
			}
			if len(merged) > 0 {
				resp.Data.Cursor = merged[len(merged)-1].PublishTime
			}
			resp.Data.IsEnd = len(pairs) <= int(in.PageSize)
			return resp, nil
		}
	}

	records, err := l.svcCtx.ConcernedRecordModel.FindByUserId(l.ctx, in.UserId, in.BizId, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(records) > int(in.PageSize) {
		records = records[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(records) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	// 组装关注动态记录列表数据
	items := make([]*social.ConcernedItem, 0, len(records))
	for _, v := range records {
		items = append(items, &social.ConcernedItem{
			Id:         v.ID,
			BizId:      v.BizID,
			ObjId:      v.ObjID,
			CreateTime: v.CreateTime.Unix(),
		})
	}

	resp.Data.Items = items
	resp.Data.Cursor = records[len(records)-1].ID
	resp.Data.IsEnd = isEnd
	return resp, nil
}
