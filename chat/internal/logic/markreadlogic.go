package logic

import (
	"context"

	"rpc-social/chat/code"
	"rpc-social/chat/internal/svc"
	"rpc-social/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MarkReadLogic) MarkRead(in *pb.MarkReadRequest) (*pb.MarkReadResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	_ = l.svcCtx.ConversationModel.ClearUnread(l.ctx, in.UserId, in.ConversationId)
	_ = l.svcCtx.MessageModel.MarkReadByConversation(l.ctx, in.ConversationId, in.UserId)

	return &pb.MarkReadResponse{}, nil
}
