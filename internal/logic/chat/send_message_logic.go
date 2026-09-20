package chatlogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/chat"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(in *social.SendMessageRequest) (resp *social.SendMessageResponse, err error) {
	resp = new(social.SendMessageResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.SendMessageData)

	if in.SenderId == 0 {
		return nil, code.SenderIdEmpty
	}
	if in.ReceiverId == 0 {
		return nil, code.ReceiverIdEmpty
	}
	if in.Content == "" {
		return nil, code.ContentEmpty
	}
	if in.SenderId == in.ReceiverId {
		return nil, code.CannotSelfChat
	}

	msg := &types.ChatMsg{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Content:    in.Content,
		MsgType:    in.MsgType,
	}

	if l.svcCtx.KqPusherClient != nil {
		threading.GoSafe(func() {
			data, err := json.Marshal(msg)
			if err != nil {
				l.Errorf("[SendMessage] marshal err: %v msg: %+v", err, msg)
				return
			}
			if err := l.svcCtx.KqPusherClient.Push(context.Background(), string(data)); err != nil {
				l.Errorf("[SendMessage] kq push err: %v data: %s", err, data)
			}
		})
	}

	return resp, nil
}
