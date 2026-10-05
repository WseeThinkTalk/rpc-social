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
		resp.Code = int64(code.BizIdEmpty.Code())
		resp.Msg = code.BizIdEmpty.Message()
		return resp, nil
	}
	if in.TargetId == 0 {
		resp.Code = int64(code.TargetIdEmpty.Code())
		resp.Msg = code.TargetIdEmpty.Message()
		return resp, nil
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
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
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

	// 2. 批量拉取根评论对应的子回复（限制返回条数）
	rootIds := make([]int64, len(roots))
	for i, v := range roots {
		rootIds[i] = v.ID
	}
	subReplies, err := l.svcCtx.ReplyModel.FindTopSubRepliesByParentIDs(l.ctx, rootIds, types.MaxSubReplyCount)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	// 3. 按 parentId 分组子回复
	subMap := make(map[int64][]*social.ReplyItem)
	for _, v := range subReplies {
		item := l.toReplyItem(v)
		subMap[v.ParentID] = append(subMap[v.ParentID], item)
	}

	// 4. 组装评论列表并绑定子回复
	items := make([]*social.ReplyItem, 0, len(roots))
	for _, v := range roots {
		rootItem := l.toReplyItem(v)
		if subs, ok := subMap[v.ID]; ok {
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
