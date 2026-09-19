package replyrpc

import (
	"rpc-social/reply/internal/config"
	"rpc-social/reply/internal/server"
	"rpc-social/reply/internal/svc"
	"rpc-social/reply/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterReplyServer(grpcServer, server.NewReplyServer(ctx))
}
