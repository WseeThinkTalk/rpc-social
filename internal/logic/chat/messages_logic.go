package chatlogic

import (
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
	resp.Code = 200
	resp.Msg = "success"
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
		l.Errorf("[Messages] FindByConversationId err: %v convId: %d", err, in.ConversationId)
		return nil, err
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

	items := make([]*social.MessageItem, 0, len(msgs))
	for _, m := range msgs {
		items = append(items, &social.MessageItem{
			Id:             m.ID,
			ConversationId: m.ConversationID,
			SenderId:       m.SenderID,
			ReceiverId:     m.ReceiverID,
			Content:        m.Content,
			MsgType:        int32(m.MsgType),
			IsRead:         m.IsRead == 1,
			CreateTime:     m.CreateTime.Unix(),
		})
	}

	resp.Data.Items = items
	resp.Data.Cursor = msgs[len(msgs)-1].ID
	resp.Data.IsEnd = isEnd
	return resp, nil
}
