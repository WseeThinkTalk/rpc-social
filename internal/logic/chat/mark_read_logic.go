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

func (l *MarkReadLogic) MarkRead(in *social.ChatMarkReadRequest) (resp *social.ChatMarkReadResponse, err error) {
	resp = new(social.ChatMarkReadResponse)

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}

	if err := l.svcCtx.ConversationModel.ClearUnread(l.ctx, in.UserId, in.ConversationId); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if err := l.svcCtx.MessageModel.MarkReadByConversation(l.ctx, in.ConversationId, in.UserId); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
