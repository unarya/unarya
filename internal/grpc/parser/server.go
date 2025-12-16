package parser

import (
	"fmt"
	"log"
	"net"

	ts "github.com/unarya/unarya/internal/parser/transport/grpc"
	"github.com/unarya/unarya/lib/proto/pb/parserpb"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewServer() *GRPCServer {
	s := grpc.NewServer()
	// Register handlers
	parserpb.RegisterParserServiceServer(s, ts.NewParserHandler())
	for name, info := range s.GetServiceInfo() {
		log.Printf("[SYSTEM]: Service %s", name)
		for _, m := range info.Methods {
			log.Printf("[SYSTEM]: RPC %s", m.Name)
		}
	}
	return &GRPCServer{server: s}
}

func (p *GRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return p.server.Serve(lis)
}
