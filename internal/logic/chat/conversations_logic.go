package chatlogic

import (
	"context"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/chat"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationsLogic {
	return &ConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConversationsLogic) Conversations(in *social.ConversationsRequest) (resp *social.ConversationsResponse, err error) {
	resp = new(social.ConversationsResponse)
	resp.Data = new(social.ConversationsData)
	resp.Data.Items = make([]*social.ConversationItem, 0)

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	convs, err := l.svcCtx.ConversationModel.FindByUserId(l.ctx, in.UserId, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(convs) > int(in.PageSize) {
		convs = convs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(convs) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	items := make([]*social.ConversationItem, 0, len(convs))
	for _, c := range convs {
		items = append(items, &social.ConversationItem{
			Id:              c.ID,
			TargetUserId:    c.TargetUserID,
			LastMessage:     c.LastMessage,
			LastMessageTime: c.LastMessageTime.Unix(),
			UnreadCount:     int64(c.UnreadCount),
		})
	}

	resp.Data.Items = items
	resp.Data.Cursor = convs[len(convs)-1].LastMessageTime.Unix()
	resp.Data.IsEnd = isEnd
	return resp, nil
}
