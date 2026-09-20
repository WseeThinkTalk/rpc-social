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

func (l *MessagesLogic) Messages(in *social.MessagesRequest) (*social.MessagesResponse, error) {
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
		return &social.MessagesResponse{
			Code: 200,
			Msg:  "success",
			Data: &social.MessagesData{
				Items: []*social.MessageItem{},
				IsEnd: true,
			},
		}, nil
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

	cursor := msgs[len(msgs)-1].ID
	return &social.MessagesResponse{
		Code: 200,
		Msg:  "success",
		Data: &social.MessagesData{
			Items:  items,
			Cursor: cursor,
			IsEnd:  isEnd,
		},
	}, nil
}
