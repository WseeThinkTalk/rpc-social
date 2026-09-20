package concernedlogic

import (
	"context"
	"encoding/json"

	"rpc-social/internal/svc"
	types "rpc-social/internal/types/concerned"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
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

func (l *CancelConcernedLogic) CancelConcerned(in *social.CancelConcernedRequest) (*social.CancelConcernedResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.ObjId == 0 {
		return nil, code.ObjIdEmpty
	}
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	msg := &types.ConcernedMsg{
		BizId:  in.BizId,
		ObjId:  in.ObjId,
		UserId: in.UserId,
		OpType: types.OpTypeCancel,
	}

	if l.svcCtx.KqPusherClient != nil {
		threading.GoSafe(func() {
			data, err := json.Marshal(msg)
			if err != nil {
				l.Errorf("[CancelConcerned] marshal msg: %v error: %v", msg, err)
				return
			}
			err = l.svcCtx.KqPusherClient.Push(context.Background(), string(data))
			if err != nil {
				l.Errorf("[CancelConcerned] kq push data: %s error: %v", data, err)
			}
		})
	}

	return &social.CancelConcernedResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
