package collector

import (
	"fmt"
	"log"
	"net"

	ts "github.com/unarya/unarya/internal/collector/transport/grpc"
	pb "github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewServer() *GRPCServer {
	s := grpc.NewServer()

	// Register handlers
	pb.RegisterCollectorServiceServer(s, ts.NewCollectorHandler())

	for name, info := range s.GetServiceInfo() {
		log.Printf("[SYSTEM]: Service: %s", name)
		for _, m := range info.Methods {
			log.Printf("[SYSTEM] RPC: %s", m.Name)
		}
	}
	return &GRPCServer{server: s}
}

func (s *GRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	return s.server.Serve(lis)
}
