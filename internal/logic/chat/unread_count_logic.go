package chatlogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadCountLogic {
	return &UnreadCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnreadCountLogic) UnreadCount(in *social.ChatUnreadCountRequest) (resp *social.ChatUnreadCountResponse, err error) {
	resp = new(social.ChatUnreadCountResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.ChatUnreadCountData)

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}

	total, err := l.svcCtx.ConversationModel.CountUnread(l.ctx, in.UserId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Data.Total = total
	return resp, nil
}
