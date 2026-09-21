package messagelogic

import (
	"context"
	"math"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/message"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationListLogic {
	return &NotificationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationListLogic) NotificationList(in *social.NotificationListRequest) (resp *social.NotificationListResponse, err error) {
	resp = new(social.NotificationListResponse)
	resp.Data = new(social.NotificationListData)
	resp.Data.Items = make([]*social.NotificationItem, 0)

	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	notifs, err := l.svcCtx.NotificationModel.FindByUserId(l.ctx, in.UserId, in.Type, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(notifs) > int(in.PageSize) {
		notifs = notifs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(notifs) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	// 组装消息通知列表数据
	items := make([]*social.NotificationItem, 0, len(notifs))
	for _, v := range notifs {
		items = append(items, &social.NotificationItem{
			Id:            v.ID,
			Type:          v.Type,
			Title:         v.Title,
			Content:       v.Content,
			RefId:         v.RefID,
			BizId:         v.BizID,
			TriggerUserId: v.TriggerUserID,
			IsRead:        v.IsRead == 1,
			CreateTime:    v.CreateTime.Unix(),
		})
	}

	resp.Data.Items = items
	resp.Data.Cursor = notifs[len(notifs)-1].ID
	resp.Data.IsEnd = isEnd
	return resp, nil
}
