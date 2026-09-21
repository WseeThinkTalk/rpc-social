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
	resp.Data = new(social.ChatUnreadCountData)

	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	total, err := l.svcCtx.ConversationModel.CountUnread(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("[UnreadCount] CountUnread err: %v userId: %d", err, in.UserId)
		return nil, err
	}

	resp.Data.Total = total
	return resp, nil
}
