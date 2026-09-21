package replylogic

import (
	"context"
	"math"

	model "rpc-social/internal/model/reply"
	"rpc-social/internal/svc"
	types "rpc-social/internal/types/reply"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyListLogic {
	return &ReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyListLogic) ReplyList(in *social.ReplyListRequest) (resp *social.ReplyListResponse, err error) {
	resp = new(social.ReplyListResponse)
	resp.Data = new(social.ReplyListData)
	resp.Data.Items = make([]*social.ReplyItem, 0)

	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	// 1. 查询根评论
	roots, err := l.svcCtx.ReplyModel.FindRootReplies(l.ctx, in.BizId, in.TargetId, int(in.SortType), in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[ReplyList] FindRootReplies err: %v req: %+v", err, in)
		return nil, err
	}

	var (
		isEnd  bool
		cursor int64
	)
	if len(roots) > int(in.PageSize) {
		roots = roots[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(roots) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	// 2. 收集根评论 ID, 查询子回复
	rootIds := make([]int64, len(roots))
	for i, r := range roots {
		rootIds[i] = r.ID
	}
	subReplies, err := l.svcCtx.ReplyModel.FindByParentIDs(l.ctx, rootIds)
	if err != nil {
		l.Errorf("[ReplyList] FindByParentIDs err: %v rootIds: %v", err, rootIds)
	}

	// 3. 构建 parentId → subReplies 映射
	subMap := make(map[int64][]*social.ReplyItem)
	for _, sub := range subReplies {
		item := l.toReplyItem(sub)
		subMap[sub.ParentID] = append(subMap[sub.ParentID], item)
	}

	// 4. 组装返回
	items := make([]*social.ReplyItem, 0, len(roots))
	for _, root := range roots {
		rootItem := l.toReplyItem(root)
		if subs, ok := subMap[root.ID]; ok {
			if len(subs) > types.MaxSubReplyCount {
				subs = subs[:types.MaxSubReplyCount]
			}
			rootItem.SubReplies = subs
		}
		items = append(items, rootItem)
	}

	if in.SortType == types.SortByLike {
		last := roots[len(roots)-1]
		cursor = int64(last.LikeNum)
	} else {
		cursor = roots[len(roots)-1].ID
	}

	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}

func (l *ReplyListLogic) toReplyItem(r *model.Reply) *social.ReplyItem {
	return &social.ReplyItem{
		ReplyId:       r.ID,
		BizId:         r.BizID,
		TargetId:      r.TargetID,
		ReplyUserId:   r.ReplyUserID,
		BeReplyUserId: r.BeReplyUserID,
		ParentId:      r.ParentID,
		Content:       r.Content,
		LikeNum:       int64(r.LikeNum),
		CreateTime:    r.CreateTime.Unix(),
	}
}
