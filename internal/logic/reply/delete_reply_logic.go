package replylogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/reply"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteReplyLogic {
	return &DeleteReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteReplyLogic) DeleteReply(in *social.DeleteReplyRequest) (resp *social.DeleteReplyResponse, err error) {
	resp = new(social.DeleteReplyResponse)

	if in.ReplyId == 0 {
		resp.Code = int64(code.ReplyNotFound.Code())
		resp.Msg = code.ReplyNotFound.Message()
		return resp, nil
	}
	if in.UserId == 0 {
		resp.Code = int64(code.CannotDeleteReply.Code())
		resp.Msg = code.CannotDeleteReply.Message()
		return resp, nil
	}

	reply, err := l.svcCtx.ReplyModel.FindOne(l.ctx, in.ReplyId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if reply == nil {
		resp.Code = int64(code.ReplyNotFound.Code())
		resp.Msg = code.ReplyNotFound.Message()
		return resp, nil
	}
	if !in.IsAdmin && reply.ReplyUserID != in.UserId {
		resp.Code = int64(code.CannotDeleteReply.Code())
		resp.Msg = code.CannotDeleteReply.Message()
		return resp, nil
	}

	msg := &types.ReplyMsg{
		ReplyId: in.ReplyId,
		UserId:  reply.ReplyUserID,
		OpType:  types.OpTypeDelete,
	}

	if l.svcCtx.KqPusherClient != nil {
		data, err := json.Marshal(msg)
		if err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			return resp, nil
		}
		if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(data)); err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			return resp, nil
		}
	}

	return resp, nil
}
