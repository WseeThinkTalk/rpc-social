package main

import (
	"flag"
	"fmt"

	chatrpc "rpc-social/chat"
	concernedrpc "rpc-social/concerned"
	likerpc "rpc-social/like"
	messagerpc "rpc-social/message"
	"rpc-social/pkg/env"
	"rpc-social/pkg/interceptors"
	replyrpc "rpc-social/reply"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/social.yaml", "the config file")

type Config struct {
	zrpc.RpcServerConf
	DataSource   string
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	Mysql sqlx.SqlConf
	DB    struct {
		DataSource   string
		MaxOpenConns int `json:",default=10"`
		MaxIdleConns int `json:",default=100"`
		MaxLifetime  int `json:",default=3600"`
	}
	CacheRedis cache.CacheConf
	BizRedis   redis.RedisConf
}

func main() {
	flag.Parse()

	env.LoadEnv()

	var c Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if c.DB.DataSource == "" {
		c.DB.DataSource = c.DataSource
	}
	if c.Mysql.DataSource == "" {
		c.Mysql.DataSource = c.DataSource
	}

	likeConf := likerpc.Config{
		RpcServerConf: c.RpcServerConf,
		KqPusherConf:  c.KqPusherConf,
		Mysql:         c.Mysql,
		CacheRedis:    c.CacheRedis,
	}

	concernedConf := concernedrpc.Config{
		RpcServerConf: c.RpcServerConf,
		KqPusherConf:  c.KqPusherConf,
		DB:            c.DB,
		BizRedis:      c.BizRedis,
	}

	replyConf := replyrpc.Config{
		RpcServerConf: c.RpcServerConf,
		KqPusherConf:  c.KqPusherConf,
		DB:            c.DB,
		BizRedis:      c.BizRedis,
	}

	messageConf := messagerpc.Config{
		RpcServerConf: c.RpcServerConf,
		DB:            c.DB,
		BizRedis:      c.BizRedis,
	}

	chatConf := chatrpc.Config{
		RpcServerConf: c.RpcServerConf,
		KqPusherConf:  c.KqPusherConf,
		DB:            c.DB,
		BizRedis:      c.BizRedis,
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		likerpc.Register(grpcServer, likeConf)
		concernedrpc.Register(grpcServer, concernedConf)
		replyrpc.Register(grpcServer, replyConf)
		messagerpc.Register(grpcServer, messageConf)
		chatrpc.Register(grpcServer, chatConf)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())
	defer s.Stop()

	fmt.Printf("Starting unified social rpc server (like, concerned, reply, message, chat) at %s...\n", c.ListenOn)
	s.Start()
}
