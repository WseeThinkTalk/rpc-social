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
	resp.Data = new(social.ReplyCountData)

	if in.BizId == "" {
		resp.Code = int64(code.BizIdEmpty.Code())
		resp.Msg = code.BizIdEmpty.Message()
		return resp, nil
	}
	if in.TargetId == 0 {
		resp.Code = int64(code.TargetIdEmpty.Code())
		resp.Msg = code.TargetIdEmpty.Message()
		return resp, nil
	}

	total, err := l.svcCtx.ReplyModel.CountByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	rootTotal, err := l.svcCtx.ReplyModel.CountRootByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Data.ReplyNum = total
	resp.Data.ReplyRootNum = rootTotal
	return resp, nil
}
