package concernedlogic

import (
	"context"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/concerned"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsConcernedLogic {
	return &IsConcernedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsConcernedLogic) IsConcerned(in *social.IsConcernedRequest) (resp *social.IsConcernedResponse, err error) {
	resp = new(social.IsConcernedResponse)
	resp.Data = new(social.IsConcernedData)

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
	if in.UserId == 0 {
		resp.Code = int64(code.UserIdEmpty.Code())
		resp.Msg = code.UserIdEmpty.Message()
		return resp, nil
	}

	record, err := l.svcCtx.ConcernedRecordModel.FindByBizIDObjIDUserID(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Data.IsConcerned = record != nil && record.Status == types.StatusConcerned
	return resp, nil
}
