package main

import (
	"log"
	"os"

	"github.com/unarya/unarya/lib/utils"
	pkg "github.com/unarya/unarya/pkg/collector"
)

func main() {
	port := os.Getenv("COLLECTOR_PORT")
	if port == "" {
		port = "50051"
	}
	defer os.Exit(1)

	pkg.InitGRPCServer()
	pkg.MinIOConnect()

	log.Printf("Collector services running on port %s", port)

	if err := pkg.StartGRPC(utils.ToInt(port)); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
