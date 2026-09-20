package replylogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyCountLogic {
	return &ReplyCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyCountLogic) ReplyCount(in *social.ReplyCountRequest) (resp *social.ReplyCountResponse, err error) {
	resp = new(social.ReplyCountResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(social.ReplyCountData)

	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}

	total, err := l.svcCtx.ReplyModel.CountByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		l.Errorf("[ReplyCount] CountByBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}

	rootTotal, err := l.svcCtx.ReplyModel.CountRootByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		l.Errorf("[ReplyCount] CountRootByBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}

	resp.Data.ReplyNum = total
	resp.Data.ReplyRootNum = rootTotal
	return resp, nil
}
