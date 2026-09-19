package chatrpc

import (
	"rpc-social/chat/internal/config"
	"rpc-social/chat/internal/server"
	"rpc-social/chat/internal/svc"
	"rpc-social/chat/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterChatServer(grpcServer, server.NewChatServer(ctx))
}
