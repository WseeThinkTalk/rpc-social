package concernedlogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedCountLogic {
	return &ConcernedCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConcernedCountLogic) ConcernedCount(in *social.ConcernedCountRequest) (resp *social.ConcernedCountResponse, err error) {
	resp = new(social.ConcernedCountResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.ConcernedCountData)

	if in.BizId == "" {
		resp.Code = int64(code.BizIdEmpty.Code())
		resp.Msg = code.BizIdEmpty.Message()
		return resp, nil
	}
	if in.ObjId == 0 {
		resp.Code = int64(code.ObjIdEmpty.Code())
		resp.Msg = code.ObjIdEmpty.Message()
		return resp, nil
	}

	if in.ObjId < 0 {
		userId := -in.ObjId
		count, err := l.svcCtx.ConcernedRecordModel.CountByUserId(l.ctx, userId, in.BizId)
		if err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			return resp, nil
		}
		resp.Data.ConcernedNum = count
		return resp, nil
	}

	count, err := l.svcCtx.ConcernedCountModel.FindByBizIDAndObjID(l.ctx, in.BizId, in.ObjId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	var num int64
	if count != nil {
		num = int64(count.ConcernedNum)
	}

	resp.Data.ConcernedNum = num
	return resp, nil
}
