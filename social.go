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
	"rpc-social/pkg/lib/etcdx"
	"rpc-social/pkg/lib/zapx"
	"rpc-social/social"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runRemoteConfig() *config.Config {
	var c config.Config
	etcdx.MustLoadRemoteConfig("/thinktalk/config/social.rpc", &c)
	if c.DB.DataSource == "" {
		c.DB.DataSource = c.DataSource
	}
	if c.Mysql.DataSource == "" {
		c.Mysql.DataSource = c.DataSource
	}
	return &c
}

func main() {
	flag.Parse()

	// 从 Etcd 配置中心拉取远程配置 (Fail-Fast)
	c := runRemoteConfig()
	if c == nil {
		return
	}

	// init logger
	writer, err := zapx.NewZapWriter()
	if err == nil {
		logx.SetWriter(writer)
	}

	ctx := svc.NewServiceContext(*c)

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
	likeSrv := likeserver.NewLikeServer(ctx)
	concernedSrv := concernedserver.NewConcernedServer(ctx)
	replySrv := replyserver.NewReplyServer(ctx)
	messageSrv := messageserver.NewMessageServer(ctx)
	chatSrv := chatserver.NewChatServer(ctx)

	social.RegisterLikeServer(grpcServer, likeSrv)
	social.RegisterConcernedServer(grpcServer, concernedSrv)
	social.RegisterReplyServer(grpcServer, replySrv)
	social.RegisterMessageServer(grpcServer, messageSrv)
	social.RegisterChatServer(grpcServer, chatSrv)

	// 兼容旧版客户端 (api-thinktalk) 调用的服务命名空间
	likeDesc := social.Like_ServiceDesc
	likeDesc.ServiceName = "service.Like"
	grpcServer.RegisterService(&likeDesc, likeSrv)

	concernedDesc := social.Concerned_ServiceDesc
	concernedDesc.ServiceName = "service.Concerned"
	grpcServer.RegisterService(&concernedDesc, concernedSrv)

	replyDesc := social.Reply_ServiceDesc
	replyDesc.ServiceName = "service.Reply"
	grpcServer.RegisterService(&replyDesc, replySrv)

	messageDesc := social.Message_ServiceDesc
	messageDesc.ServiceName = "service.Message"
	grpcServer.RegisterService(&messageDesc, messageSrv)

	chatDesc := social.Chat_ServiceDesc
	chatDesc.ServiceName = "service.Chat"
	grpcServer.RegisterService(&chatDesc, chatSrv)
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

