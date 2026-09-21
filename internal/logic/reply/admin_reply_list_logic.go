package replylogic

import (
	"rpc-social/pkg/code"
	"context"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/reply"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminReplyListLogic {
	return &AdminReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminReplyListLogic) AdminReplyList(in *social.AdminReplyListRequest) (resp *social.AdminReplyListResponse, err error) {
	resp = new(social.AdminReplyListResponse)
	resp.Data = new(social.ReplyListData)
	resp.Data.Items = make([]*social.ReplyItem, 0)

	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	replies, err := l.svcCtx.ReplyModel.AdminFindAll(l.ctx, in.Keyword, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(replies) > int(in.PageSize) {
		replies = replies[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(replies) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	items := make([]*social.ReplyItem, 0, len(replies))
	for _, r := range replies {
		items = append(items, &social.ReplyItem{
			ReplyId:       r.ID,
			BizId:         r.BizID,
			TargetId:      r.TargetID,
			ReplyUserId:   r.ReplyUserID,
			BeReplyUserId: r.BeReplyUserID,
			ParentId:      r.ParentID,
			Content:       r.Content,
			LikeNum:       int64(r.LikeNum),
			CreateTime:    r.CreateTime.Unix(),
		})
	}

	var cursor int64
	if len(items) > 0 && !isEnd {
		cursor = items[len(items)-1].ReplyId
	}

	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}
