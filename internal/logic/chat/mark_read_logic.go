package chatlogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

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

func (l *MarkReadLogic) MarkRead(in *social.ChatMarkReadRequest) (*social.ChatMarkReadResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	_ = l.svcCtx.ConversationModel.ClearUnread(l.ctx, in.UserId, in.ConversationId)
	_ = l.svcCtx.MessageModel.MarkReadByConversation(l.ctx, in.ConversationId, in.UserId)

	return &social.ChatMarkReadResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
