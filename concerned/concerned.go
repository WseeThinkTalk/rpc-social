package concernedrpc

import (
	"rpc-social/concerned/internal/config"
	"rpc-social/concerned/internal/server"
	"rpc-social/concerned/internal/svc"
	"rpc-social/concerned/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterConcernedServer(grpcServer, server.NewConcernedServer(ctx))
}
