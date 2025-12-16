package collector

import (
	"github.com/unarya/unarya/internal/grpc/collector"
	"github.com/unarya/unarya/internal/infrastructures"
)

var grpcServer *collector.GRPCServer

func InitGRPCServer() {
	grpcServer = collector.NewServer()
}

func MinIOConnect() {
	infrastructures.ConnectMinio()
}

func StartGRPC(port int) error {
	return grpcServer.Start(port)
}
