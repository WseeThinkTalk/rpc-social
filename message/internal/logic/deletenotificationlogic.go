package logic

import (
	"context"

	"rpc-social/message/code"
	"rpc-social/message/internal/svc"
	"rpc-social/message/pb"

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

func (l *DeleteNotificationLogic) DeleteNotification(in *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.NotificationId == 0 {
		return nil, code.NotificationIdEmpty
	}

	// 先查通知是否存在且属于该用户
	notif, err := l.svcCtx.NotificationModel.FindOne(l.ctx, in.NotificationId)
	if err != nil {
		l.Errorf("[DeleteNotification] FindOne err: %v notificationId: %d", err, in.NotificationId)
		return nil, err
	}
	if notif == nil {
		return nil, code.NotificationNotFound
	}
	if notif.UserID != in.UserId {
		return nil, code.NotificationNotFound
	}

	if err := l.svcCtx.NotificationModel.Delete(l.ctx, in.NotificationId); err != nil {
		l.Errorf("[DeleteNotification] Delete err: %v notificationId: %d", err, in.NotificationId)
		return nil, err
	}

	return &pb.DeleteNotificationResponse{}, nil
}
