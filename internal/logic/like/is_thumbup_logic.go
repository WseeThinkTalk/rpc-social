package likelogic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "rpc-social/internal/model/like"
	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
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

	// 优先查询缓存
	cacheKey := fmt.Sprintf("biz#like#status:%s:%d:%d", in.BizId, in.TargetId, in.UserId)
	if l.svcCtx.BizRedis != nil {
		val, rerr := l.svcCtx.BizRedis.GetCtx(l.ctx, cacheKey)
		if rerr == nil && val != "" {
			likeType, _ := strconv.Atoi(val)
			if likeType > 0 {
				resp.Data.UserThumbups[in.TargetId] = &social.UserThumbup{
					UserId:      in.UserId,
					ThumbupTime: time.Now().UnixMilli(),
					LikeType:    int32(likeType),
				}
			}
			return resp, nil
		}
	}

	record, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(l.ctx, in.BizId, in.TargetId, in.UserId)
	if err != nil && err != model.ErrNotFound {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	if record != nil {
		resp.Data.UserThumbups[in.TargetId] = &social.UserThumbup{
			UserId:      record.UserId,
			ThumbupTime: record.CreateTime.UnixMilli(),
			LikeType:    int32(record.LikeType),
		}
		if l.svcCtx.BizRedis != nil {
			_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, cacheKey, strconv.Itoa(int(record.LikeType)), 86400*7)
		}
	} else if l.svcCtx.BizRedis != nil {
		// 写入空值防缓存穿透
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, cacheKey, "0", 300)
	}

	return resp, nil
}
