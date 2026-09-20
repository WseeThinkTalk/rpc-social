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

type AddConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConcernedLogic {
	return &AddConcernedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddConcernedLogic) AddConcerned(in *social.AddConcernedRequest) (resp *social.AddConcernedResponse, err error) {
	resp = new(social.AddConcernedResponse)
	resp.Code = 200
	resp.Msg = "success"

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
		OpType: types.OpTypeAdd,
	}

	if l.svcCtx.KqPusherClient != nil {
		threading.GoSafe(func() {
			data, err := json.Marshal(msg)
			if err != nil {
				l.Errorf("[AddConcerned] marshal msg: %v error: %v", msg, err)
				return
			}
			err = l.svcCtx.KqPusherClient.Push(context.Background(), string(data))
			if err != nil {
				l.Errorf("[AddConcerned] kq push data: %s error: %v", data, err)
			}
		})
	}

	return resp, nil
}
