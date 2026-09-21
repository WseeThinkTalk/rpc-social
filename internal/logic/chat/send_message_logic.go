package chatlogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/chat"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
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
	resp.Data = new(social.SendMessageData)

	if in.SenderId == 0 {
		resp.Code = int64(code.SenderIdEmpty.Code())
		resp.Msg = code.SenderIdEmpty.Message()
		return resp, nil
	}
	if in.ReceiverId == 0 {
		resp.Code = int64(code.ReceiverIdEmpty.Code())
		resp.Msg = code.ReceiverIdEmpty.Message()
		return resp, nil
	}
	if in.Content == "" {
		resp.Code = int64(code.ContentEmpty.Code())
		resp.Msg = code.ContentEmpty.Message()
		return resp, nil
	}
	if in.SenderId == in.ReceiverId {
		resp.Code = int64(code.CannotSelfChat.Code())
		resp.Msg = code.CannotSelfChat.Message()
		return resp, nil
	}

	msg := &types.ChatMsg{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Content:    in.Content,
		MsgType:    in.MsgType,
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
