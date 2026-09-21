package concernedlogic

import (
	"context"
	"math"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/concerned"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedListLogic {
	return &ConcernedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConcernedListLogic) ConcernedList(in *social.ConcernedListRequest) (resp *social.ConcernedListResponse, err error) {
	resp = new(social.ConcernedListResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.ConcernedListData)
	resp.Data.Items = make([]*social.ConcernedItem, 0)

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

	records, err := l.svcCtx.ConcernedRecordModel.FindByUserId(l.ctx, in.UserId, in.BizId, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(records) > int(in.PageSize) {
		records = records[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(records) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	items := make([]*social.ConcernedItem, 0, len(records))
	for _, r := range records {
		items = append(items, &social.ConcernedItem{
			Id:         r.ID,
			BizId:      r.BizID,
			ObjId:      r.ObjID,
			CreateTime: r.CreateTime.Unix(),
		})
	}

	resp.Data.Items = items
	resp.Data.Cursor = records[len(records)-1].ID
	resp.Data.IsEnd = isEnd
	return resp, nil
}
