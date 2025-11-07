package collector

import (
	"fmt"
	"log"
	"net"

	collector2 "github.com/unarya/unarya/internal/collector/transport/grpc"
	pb "github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewServer() *GRPCServer {
	s := grpc.NewServer()

	// Register handlers
	pb.RegisterCollectorServiceServer(s, collector2.NewCollectorHandler())

	return &GRPCServer{server: s}
}

func (s *GRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}
	return s.server.Serve(lis)
}
