package likelogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/like"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *social.ThumbupRequest) (resp *social.ThumbupResponse, err error) {
	resp = new(social.ThumbupResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.ThumbupData)
	resp.Data.BizId = in.BizId
	resp.Data.ObjId = in.ObjId

	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}

	if l.svcCtx.KqPusherClient != nil {
		threading.GoSafe(func() {
			data, err := json.Marshal(msg)
			if err != nil {
				l.Errorf("[Thumbup] marshal msg: %v error: %v", msg, err)
				return
			}

			err = l.svcCtx.KqPusherClient.Push(context.Background(), string(data))
			if err != nil {
				l.Errorf("[Thumbup] kq push data: %s error: %v", data, err)
			}
		})
	}

	return resp, nil
}
