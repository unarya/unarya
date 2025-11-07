package collector

import (
	"github.com/unarya/unarya/internal/grpc/collector"
)

var grpcServer *collector.GRPCServer

func InitGRPCServer() {
	grpcServer = collector.NewServer()
}

func StartGRPC(port int) error {
	return grpcServer.Start(port)
}
