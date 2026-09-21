package chatlogic

import (
	"rpc-social/pkg/code"
	"context"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/chat"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagesLogic {
	return &MessagesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MessagesLogic) Messages(in *social.MessagesRequest) (resp *social.MessagesResponse, err error) {
	resp = new(social.MessagesResponse)
	resp.Data = new(social.MessagesData)
	resp.Data.Items = make([]*social.MessageItem, 0)

	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize
	}

	msgs, err := l.svcCtx.MessageModel.FindByConversationId(l.ctx, in.ConversationId, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(msgs) > int(in.PageSize) {
		msgs = msgs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(msgs) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	// 组装私信聊天消息列表数据
	items := make([]*social.MessageItem, 0, len(msgs))
	for _, v := range msgs {
		items = append(items, &social.MessageItem{
			Id:             v.ID,
			ConversationId: v.ConversationID,
			SenderId:       v.SenderID,
			ReceiverId:     v.ReceiverID,
			Content:        v.Content,
			MsgType:        int32(v.MsgType),
			IsRead:         v.IsRead == 1,
			CreateTime:     v.CreateTime.Unix(),
		})
	}

	resp.Data.Items = items
	resp.Data.Cursor = msgs[len(msgs)-1].ID
	resp.Data.IsEnd = isEnd
	return resp, nil
}
