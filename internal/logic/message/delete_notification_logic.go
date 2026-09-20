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

func (l *DeleteNotificationLogic) DeleteNotification(in *social.DeleteNotificationRequest) (*social.DeleteNotificationResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.NotificationId == 0 {
		return nil, code.NotificationIdEmpty
	}

	notif, err := l.svcCtx.NotificationModel.FindOne(l.ctx, in.NotificationId)
	if err != nil {
		l.Errorf("[DeleteNotification] FindOne err: %v notificationId: %d", err, in.NotificationId)
		return nil, err
	}
	if notif == nil || notif.UserID != in.UserId {
		return nil, code.NotificationNotFound
	}

	if err := l.svcCtx.NotificationModel.Delete(l.ctx, in.NotificationId); err != nil {
		l.Errorf("[DeleteNotification] Delete err: %v notificationId: %d", err, in.NotificationId)
		return nil, err
	}

	return &social.DeleteNotificationResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
