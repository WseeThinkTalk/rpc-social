package svc

import (
	"rpc-social/like/internal/config"
	"rpc-social/like/internal/model"
	"rpc-social/pkg/orm"

	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config          config.Config
	KqPusherClient  *kq.Pusher
	LikeRecordModel model.LikeRecordModel
	DB              *orm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewPostgres(&orm.Config{
		DSN:          c.Mysql.DataSource,
		MaxOpenConns: 10,
		MaxIdleConns: 100,
		MaxLifetime:  3600,
	})
	return &ServiceContext{
		Config:          c,
		KqPusherClient:  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		LikeRecordModel: model.NewLikeRecordModel(db.DB),
		DB:              db,
	}
}
