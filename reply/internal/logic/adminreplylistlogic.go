package logic

import (
	"context"

	"rpc-social/reply/internal/svc"
	"rpc-social/reply/internal/types"
	"rpc-social/reply/pb"

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

func (l *AdminReplyListLogic) AdminReplyList(in *pb.AdminReplyListRequest) (*pb.AdminReplyListResponse, error) {
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	replies, err := l.svcCtx.ReplyModel.AdminFindAll(l.ctx, in.Keyword, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[AdminReplyList] AdminFindAll err: %v req: %+v", err, in)
		return nil, err
	}

	var isEnd bool
	if len(replies) > int(in.PageSize) {
		replies = replies[:in.PageSize]
	} else {
		isEnd = true
	}

	items := make([]*pb.ReplyItem, 0, len(replies))
	for _, r := range replies {
		items = append(items, &pb.ReplyItem{
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

	return &pb.AdminReplyListResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}
