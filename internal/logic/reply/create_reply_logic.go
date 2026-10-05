package replylogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/reply"
	"rpc-social/pkg/code"
	"rpc-social/pkg/guard"
	"rpc-social/pkg/sensitive"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var replyGuard = guard.NewMemoryGuard(2 * time.Second)

type CreateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateReplyLogic) CreateReply(in *social.CreateReplyRequest) (resp *social.CreateReplyResponse, err error) {
	resp = new(social.CreateReplyResponse)
	resp.Data = new(social.CreateReplyData)

	if in.BizId == "" {
		resp.Code = int64(code.BizIdEmpty.Code())
		resp.Msg = code.BizIdEmpty.Message()
		return resp, nil
	}
	if in.TargetId == 0 {
		resp.Code = int64(code.TargetIdEmpty.Code())
		resp.Msg = code.TargetIdEmpty.Message()
		return resp, nil
	}
	if in.ReplyUserId == 0 {
		resp.Code = int64(code.ReplyUserIdEmpty.Code())
		resp.Msg = code.ReplyUserIdEmpty.Message()
		return resp, nil
	}
	if in.Content == "" {
		resp.Code = int64(code.ContentEmpty.Code())
		resp.Msg = code.ContentEmpty.Message()
		return resp, nil
	}
	if len(in.Content) > 5000 {
		resp.Code = int64(code.ContentTooLong.Code())
		resp.Msg = code.ContentTooLong.Message()
		return resp, nil
	}

	// 频控检查（同一用户 2 秒内限发 1 条评论）
	guardKey := fmt.Sprintf("%d:create_reply", in.ReplyUserId)
	if !replyGuard.Acquire(guardKey) {
		resp.Code = int64(code.FrequentOperation.Code())
		resp.Msg = code.FrequentOperation.Message()
		return resp, nil
	}

	// 敏感词检查
	defaultFilter := sensitive.NewFilter([]string{"涉黄", "涉暴", "赌博", "违禁"})
	if defaultFilter.IsSensitive(in.Content) {
		resp.Code = int64(code.ReplyContainsSensitiveWord.Code())
		resp.Msg = code.ReplyContainsSensitiveWord.Message()
		return resp, nil
	}

	msg := &types.ReplyMsg{
		BizId:         in.BizId,
		TargetId:      in.TargetId,
		ReplyUserId:   in.ReplyUserId,
		BeReplyUserId: in.BeReplyUserId,
		ParentId:      in.ParentId,
		Content:       in.Content,
		OpType:        types.OpTypeCreate,
	}

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
