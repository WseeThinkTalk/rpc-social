package replylogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/reply"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
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
		return nil, code.ReplyNotFound
	}
	if in.UserId == 0 {
		return nil, code.CannotDeleteReply
	}

	reply, err := l.svcCtx.ReplyModel.FindOne(l.ctx, in.ReplyId)
	if err != nil {
		l.Errorf("[DeleteReply] ReplyModel.FindOne err: %v replyId: %d", err, in.ReplyId)
		return nil, err
	}
	if reply == nil {
		return nil, code.ReplyNotFound
	}
	if !in.IsAdmin && reply.ReplyUserID != in.UserId {
		return nil, code.CannotDeleteReply
	}

	msg := &types.ReplyMsg{
		ReplyId: in.ReplyId,
		UserId:  reply.ReplyUserID,
		OpType:  types.OpTypeDelete,
	}

	if l.svcCtx.KqPusherClient != nil {
		threading.GoSafe(func() {
			data, err := json.Marshal(msg)
			if err != nil {
				l.Errorf("[DeleteReply] marshal msg: %v error: %v", msg, err)
				return
			}
			err = l.svcCtx.KqPusherClient.Push(context.Background(), string(data))
			if err != nil {
				l.Errorf("[DeleteReply] kq push data: %s error: %v", data, err)
			}
		})
	}

	return resp, nil
}
