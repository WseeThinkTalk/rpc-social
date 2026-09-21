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

func (l *IsThumbupLogic) IsThumbup(in *social.IsThumbupRequest) (resp *social.IsThumbupResponse, err error) {
	resp = new(social.IsThumbupResponse)
	resp.Data = new(social.IsThumbupData)
	resp.Data.UserThumbups = make(map[int64]*social.UserThumbup)

	record, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(l.ctx, in.BizId, in.TargetId, in.UserId)
	if err != nil && err != model.ErrNotFound {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	if record != nil {
		resp.Data.UserThumbups[in.TargetId] = &social.UserThumbup{
			UserId:      record.UserId,
			ThumbupTime: record.CreateTime.UnixMilli(),
			LikeType:    int32(record.LikeType),
		}
	}

	return resp, nil
}
