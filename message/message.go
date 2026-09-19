package messagerpc

import (
	"rpc-social/message/internal/config"
	"rpc-social/message/internal/server"
	"rpc-social/message/internal/svc"
	"rpc-social/message/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterMessageServer(grpcServer, server.NewMessageServer(ctx))
}
