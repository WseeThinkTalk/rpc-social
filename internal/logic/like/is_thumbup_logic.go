package likelogic

import (
	"context"

	model "rpc-social/internal/model/like"
	"rpc-social/internal/svc"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsThumbupLogic) IsThumbup(in *social.IsThumbupRequest) (*social.IsThumbupResponse, error) {
	record, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(l.ctx, in.BizId, in.TargetId, in.UserId)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("[IsThumbup] find like record error: %v", err)
		return nil, err
	}

	userThumbups := make(map[int64]*social.UserThumbup)
	if record != nil {
		userThumbups[in.TargetId] = &social.UserThumbup{
			UserId:      record.UserId,
			ThumbupTime: record.CreateTime.UnixMilli(),
			LikeType:    int32(record.LikeType),
		}
	}

	return &social.IsThumbupResponse{
		Code: 200,
		Msg:  "success",
		Data: &social.IsThumbupData{
			UserThumbups: userThumbups,
		},
	}, nil
}
