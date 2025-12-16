package parser

import (
	"github.com/unarya/unarya/internal/grpc/parser"
	"github.com/unarya/unarya/internal/infrastructures"
)

var grpcServer *parser.GRPCServer

func InitGRPCServer() {
	grpcServer = parser.NewServer()
}

func StartGRPCServer(port int) error {
	return grpcServer.Start(port)
}

func MinIOConnect() {
	infrastructures.ConnectMinio()
}
