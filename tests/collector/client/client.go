package client

import (
	"fmt"
	"log"
	"os"

	"github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CollectorClient wraps the gRPC client with advanced features
type CollectorClient struct {
	Client      collectorpb.CollectorServiceClient
	Conn        *grpc.ClientConn
	GithubToken string
	Workers     int
}

// NewCollectorClient creates a new collector client
func NewCollectorClient(addr, githubToken string, workers int) (*CollectorClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("[ERROR]: cannot connect to collector services: %w", err)
	}

	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}
	log.Printf("GITHUB_TOKEN = '%s'", githubToken)
	return &CollectorClient{
		Client:      collectorpb.NewCollectorServiceClient(conn),
		Conn:        conn,
		GithubToken: githubToken,
		Workers:     workers,
	}, nil
}

// Close closes the gRPC connection
func (c *CollectorClient) Close() error {
	return c.Conn.Close()
}
