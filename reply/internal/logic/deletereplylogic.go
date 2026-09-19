package logic

import (
	"context"
	"encoding/json"

	"rpc-social/reply/code"
	"rpc-social/reply/internal/svc"
	"rpc-social/reply/pb"

	"rpc-social/reply/internal/types"

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

func (l *DeleteReplyLogic) DeleteReply(in *pb.DeleteReplyRequest) (*pb.DeleteReplyResponse, error) {
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

	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Errorf("[DeleteReply] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(data))
		if err != nil {
			l.Errorf("[DeleteReply] kq push data: %s error: %v", data, err)
		}
	})

	return &pb.DeleteReplyResponse{}, nil
}
