package replylogic

import (
	"context"

	"rpc-social/internal/svc"
	"rpc-social/pkg/code"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyDetailLogic {
	return &ReplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyDetailLogic) ReplyDetail(in *social.ReplyDetailRequest) (resp *social.ReplyDetailResponse, err error) {
	resp = new(social.ReplyDetailResponse)
	resp.Data = new(social.ReplyItem)
	resp.Data.SubReplies = make([]*social.ReplyItem, 0)

	if in.ReplyId == 0 {
		resp.Code = int64(code.ReplyNotFound.Code())
		resp.Msg = code.ReplyNotFound.Message()
		return resp, nil
	}

	reply, err := l.svcCtx.ReplyModel.FindOne(l.ctx, in.ReplyId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	if reply == nil || reply.Status == 1 {
		resp.Code = int64(code.ReplyNotFound.Code())
		resp.Msg = code.ReplyNotFound.Message()
		return resp, nil
	}

	resp.Data.ReplyId = reply.ID
	resp.Data.BizId = reply.BizID
	resp.Data.TargetId = reply.TargetID
	resp.Data.ReplyUserId = reply.ReplyUserID
	resp.Data.BeReplyUserId = reply.BeReplyUserID
	resp.Data.ParentId = reply.ParentID
	resp.Data.Content = reply.Content
	resp.Data.LikeNum = int64(reply.LikeNum)
	resp.Data.CreateTime = reply.CreateTime.Unix()

	return resp, nil
}
