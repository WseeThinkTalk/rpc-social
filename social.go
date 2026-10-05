package main

import (
	"context"
	"flag"
	"fmt"

	"rpc-social/internal/config"
	chatserver "rpc-social/internal/server/chat"
	concernedserver "rpc-social/internal/server/concerned"
	likeserver "rpc-social/internal/server/like"
	messageserver "rpc-social/internal/server/message"
	replyserver "rpc-social/internal/server/reply"
	"rpc-social/internal/svc"
	"rpc-social/pkg/env"
	"rpc-social/pkg/lib/zapx"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/social.yaml", "the config file")

func main() {
	flag.Parse()

	env.LoadEnv()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if c.DB.DataSource == "" {
		c.DB.DataSource = c.DataSource
	}
	if c.Mysql.DataSource == "" {
		c.Mysql.DataSource = c.DataSource
	}

	// init logger
	writer, err := zapx.NewZapWriter()
	if err == nil {
		logx.SetWriter(writer)
	}

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		registerServer(ctx, grpcServer)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(unaryServerInterceptor())
	defer s.Stop()

	fmt.Printf("Starting unified social rpc server (like, concerned, reply, message, chat) at %s...\n", c.ListenOn)
	s.Start()
}

// registerServer 注册 RPC 服务
func registerServer(ctx *svc.ServiceContext, grpcServer grpc.ServiceRegistrar) {
	social.RegisterLikeServer(grpcServer, likeserver.NewLikeServer(ctx))
	social.RegisterConcernedServer(grpcServer, concernedserver.NewConcernedServer(ctx))
	social.RegisterReplyServer(grpcServer, replyserver.NewReplyServer(ctx))
	social.RegisterMessageServer(grpcServer, messageserver.NewMessageServer(ctx))
	social.RegisterChatServer(grpcServer, chatserver.NewChatServer(ctx))
}

// unaryServerInterceptor grpc 拦截器
func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("[RPC Panic] info: %s, recover: %v", info.FullMethod, r)
			}
		}()

		resp, err := handler(ctx, req)
		return resp, err
	}
}

