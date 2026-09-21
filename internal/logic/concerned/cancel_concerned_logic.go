package concernedlogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/concerned"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelConcernedLogic {
	return &CancelConcernedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelConcernedLogic) CancelConcerned(in *social.CancelConcernedRequest) (resp *social.CancelConcernedResponse, err error) {
	resp = new(social.CancelConcernedResponse)

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

	msg := &types.ConcernedMsg{
		BizId:  in.BizId,
		ObjId:  in.ObjId,
		UserId: in.UserId,
		OpType: types.OpTypeCancel,
	}

	if l.svcCtx.KqPusherClient != nil {
		data, err := json.Marshal(msg)
		if err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}
		if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(data)); err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}
	}

	return resp, nil
}
