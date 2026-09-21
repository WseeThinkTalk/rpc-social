package messagelogic

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
	return &UnreadCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnreadCountLogic) UnreadCount(in *social.UnreadCountRequest) (resp *social.UnreadCountResponse, err error) {
	resp = new(social.UnreadCountResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.UnreadCountData)

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}

	total, err := l.svcCtx.NotificationModel.CountUnread(l.ctx, in.UserId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	typeCounts, err := l.svcCtx.NotificationModel.CountUnreadByType(l.ctx, in.UserId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Data.Total = total
	resp.Data.TypeCounts = typeCounts
	return resp, nil
}
