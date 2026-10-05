package likelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/like"
	"rpc-social/pkg/code"
	"rpc-social/pkg/guard"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var likeGuard = guard.NewMemoryGuard(1 * time.Second)

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

	// 频控与防连击拦截
	guardKey := fmt.Sprintf("%d:%s:%d", in.UserId, in.BizId, in.ObjId)
	if !likeGuard.Acquire(guardKey) {
		resp.Code = int64(code.FrequentOperation.Code())
		resp.Msg = code.FrequentOperation.Message()
		return resp, nil
	}

	// 构造点赞消息
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}

	// 更新缓存状态
	if l.svcCtx.BizRedis != nil {
		key := fmt.Sprintf("biz#like#status:%s:%d:%d", in.BizId, in.ObjId, in.UserId)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, strconv.Itoa(int(in.LikeType)), 86400*7)
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
