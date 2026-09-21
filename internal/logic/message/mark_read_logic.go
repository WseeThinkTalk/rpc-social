package messagelogic

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
	return &MarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkReadLogic) MarkRead(in *social.MarkReadRequest) (resp *social.MarkReadResponse, err error) {
	resp = new(social.MarkReadResponse)

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
	if notif == nil {
		resp.Code = int64(code.NotificationNotFound.Code())
		resp.Msg = code.NotificationNotFound.Message()
		return resp, nil
	}

	err = l.svcCtx.NotificationModel.UpdateRead(l.ctx, in.NotificationId, in.UserId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
