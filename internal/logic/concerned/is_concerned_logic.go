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

func (l *IsConcernedLogic) IsConcerned(in *social.IsConcernedRequest) (*social.IsConcernedResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.ObjId == 0 {
		return nil, code.ObjIdEmpty
	}
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	record, err := l.svcCtx.ConcernedRecordModel.FindByBizIDObjIDUserID(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		l.Errorf("[IsConcerned] FindByBizIDObjIDUserID err: %v req: %+v", err, in)
		return nil, err
	}

	isConcerned := record != nil && record.Status == types.StatusConcerned
	return &social.IsConcernedResponse{
		Code: 200,
		Msg:  "success",
		Data: &social.IsConcernedData{
			IsConcerned: isConcerned,
		},
	}, nil
}
