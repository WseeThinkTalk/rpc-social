package svc

import (
	"rpc-social/internal/config"
	chatmodel "rpc-social/internal/model/chat"
	concernedmodel "rpc-social/internal/model/concerned"
	likemodel "rpc-social/internal/model/like"
	messagemodel "rpc-social/internal/model/message"
	replymodel "rpc-social/internal/model/reply"
	"rpc-social/pkg/orm"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config         config.Config
	DB             *orm.DB
	BizRedis       *redis.Redis
	KqPusherClient *kq.Pusher

	// Like
	LikeRecordModel likemodel.LikeRecordModel
	LikeCountModel  likemodel.LikeCountModel

	// Concerned
	ConcernedRecordModel *concernedmodel.ConcernedRecordModel
	ConcernedCountModel  *concernedmodel.ConcernedCountModel

	// Reply
	ReplyModel      *replymodel.ReplyModel
	ReplyCountModel *replymodel.ReplyCountModel

	// Message
	NotificationModel *messagemodel.NotificationModel

	// Chat
	ConversationModel *chatmodel.ConversationModel
	MessageModel      *chatmodel.MessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := c.DB.DataSource
	if dsn == "" {
		dsn = c.DataSource
	}
	if dsn == "" {
		dsn = c.Mysql.DataSource
	}

	db := orm.MustNewPostgres(&orm.Config{
		DSN:          dsn,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})

	var rds *redis.Redis
	if c.BizRedis.Host != "" {
		rds = redis.MustNewRedis(redis.RedisConf{
			Host:        c.BizRedis.Host,
			Pass:        c.BizRedis.Pass,
			Type:        c.BizRedis.Type,
			PingTimeout: 10000000000,
		})
	}

	var pusher *kq.Pusher
	if len(c.KqPusherConf.Brokers) > 0 && c.KqPusherConf.Topic != "" {
		pusher = kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic)
	}

	return &ServiceContext{
		Config:               c,
		DB:                   db,
		BizRedis:             rds,
		KqPusherClient:       pusher,
		LikeRecordModel:      likemodel.NewLikeRecordModel(db.DB),
		LikeCountModel:       likemodel.NewLikeCountModel(db.DB),
		ConcernedRecordModel: concernedmodel.NewConcernedRecordModel(db.DB),
		ConcernedCountModel:  concernedmodel.NewConcernedCountModel(db.DB),
		ReplyModel:           replymodel.NewReplyModel(db.DB),
		ReplyCountModel:      replymodel.NewReplyCountModel(db.DB),
		NotificationModel:    messagemodel.NewNotificationModel(db.DB),
		ConversationModel:    chatmodel.NewConversationModel(db.DB),
		MessageModel:         chatmodel.NewMessageModel(db.DB),
	}
}
