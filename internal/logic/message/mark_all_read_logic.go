package messagelogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkAllReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllReadLogic {
	return &MarkAllReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkAllReadLogic) MarkAllRead(in *social.MarkAllReadRequest) (resp *social.MarkAllReadResponse, err error) {
	resp = new(social.MarkAllReadResponse)

	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	err = l.svcCtx.NotificationModel.UpdateAllRead(l.ctx, in.UserId, in.Type)
	if err != nil {
		l.Errorf("[MarkAllRead] UpdateAllRead err: %v userId: %d type: %d", err, in.UserId, in.Type)
		return nil, err
	}

	return resp, nil
}
