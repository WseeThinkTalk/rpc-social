package likerpc

import (
	"rpc-social/like/internal/config"
	"rpc-social/like/internal/server"
	"rpc-social/like/internal/svc"
	"rpc-social/like/service"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	service.RegisterLikeServer(grpcServer, server.NewLikeServer(ctx))
}
