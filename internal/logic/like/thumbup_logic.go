package likelogic

import (
	"rpc-social/pkg/code"
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/like"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
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

// Thumbup 点赞或取消点赞
func (l *ThumbupLogic) Thumbup(in *social.ThumbupRequest) (resp *social.ThumbupResponse, err error) {
	resp = new(social.ThumbupResponse)
	resp.Data = new(social.ThumbupData)
	resp.Data.BizId = in.BizId
	resp.Data.ObjId = in.ObjId

	// 构造点赞消息
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}

	// 异步投递到消息队列
	if l.svcCtx.KqPusherClient != nil {
		data, err := json.Marshal(msg)
		if err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}

		if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(data)); err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}
	}

	return resp, nil
}
