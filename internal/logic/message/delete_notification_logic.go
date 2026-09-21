package messagelogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationLogic {
	return &DeleteNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNotificationLogic) DeleteNotification(in *social.DeleteNotificationRequest) (resp *social.DeleteNotificationResponse, err error) {
	resp = new(social.DeleteNotificationResponse)
	resp.Code = 200
	resp.Msg = "success"

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}
	if in.NotificationId == 0 {
		resp.Code = int64(code.NotificationIdEmpty.Code())
		resp.Msg = code.NotificationIdEmpty.Message()
		return resp, nil
	}

	notif, err := l.svcCtx.NotificationModel.FindOne(l.ctx, in.NotificationId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if notif == nil || notif.UserID != in.UserId {
		resp.Code = int64(code.NotificationNotFound.Code())
		resp.Msg = code.NotificationNotFound.Message()
		return resp, nil
	}

	if err := l.svcCtx.NotificationModel.Delete(l.ctx, in.NotificationId); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
