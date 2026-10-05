package code

import "rpc-social/pkg/xcode"

var (
	// Common
	ServerErr = xcode.ServerErr
	NotFound  = xcode.NotFound

	// Chat (10000+)
	ChatSenderIdEmpty   = xcode.New(10001, "发送者ID不能为空")
	ChatReceiverIdEmpty = xcode.New(10002, "接收者ID不能为空")
	ChatContentEmpty    = xcode.New(10003, "消息内容不能为空")
	ChatUserIdEmpty     = xcode.New(10004, "用户ID不能为空")
	CannotSelfChat      = xcode.New(10005, "不能给自己发消息")

	// Concerned (80000+)
	ConcernedBizIdEmpty  = xcode.New(80001, "业务ID不能为空")
	ConcernedObjIdEmpty  = xcode.New(80002, "收藏对象ID不能为空")
	ConcernedUserIdEmpty = xcode.New(80003, "用户ID不能为空")

	// Message (70000+)
	MessageUserIdEmpty   = xcode.New(70001, "用户ID不能为空")
	NotificationIdEmpty  = xcode.New(70002, "通知ID不能为空")
	NotificationNotFound = xcode.New(70003, "通知不存在")

	// Reply (60000+)
	ReplyBizIdEmpty     = xcode.New(60001, "业务ID不能为空")
	ReplyTargetIdEmpty  = xcode.New(60002, "评论目标ID不能为空")
	ReplyUserIdEmpty    = xcode.New(60003, "评论用户ID不能为空")
	ReplyContentEmpty   = xcode.New(60004, "评论内容不能为空")
	ReplyContentTooLong = xcode.New(60005, "评论内容过长")
	ReplyNotFound       = xcode.New(60006, "评论不存在")
	CannotDeleteReply          = xcode.New(60007, "无权删除此评论")
	ReplyContainsSensitiveWord = xcode.New(60008, "评论内容包含违规敏感词汇")

	// Common Aliases
	UserIdEmpty     = ChatUserIdEmpty
	SenderIdEmpty   = ChatSenderIdEmpty
	ReceiverIdEmpty = ChatReceiverIdEmpty
	ContentEmpty    = ChatContentEmpty
	BizIdEmpty      = ConcernedBizIdEmpty
	ObjIdEmpty      = ConcernedObjIdEmpty
	TargetIdEmpty   = ReplyTargetIdEmpty
	ContentTooLong  = ReplyContentTooLong
)
