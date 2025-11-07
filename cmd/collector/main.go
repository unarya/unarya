package main

import (
	"log"
	"os"

	"github.com/unarya/unarya/lib/utils"
	pkg "github.com/unarya/unarya/pkg/collector"
)

// main starts the gRPC Collector service
func main() {
	port := os.Getenv("COLLECTOR_PORT")
	if port == "" {
		port = "50051"
	}

	pkg.InitGRPCServer()

	if err := pkg.StartGRPC(utils.ToInt(port)); err != nil {
		log.Fatalf("Failed to start gRPC server on port %s: %v", port, err)
	}

	log.Printf("Collector service running on port %s", port)
}
