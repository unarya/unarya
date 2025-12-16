package main

import (
	"log"
	"os"

	"github.com/unarya/unarya/lib/utils"
	pkg "github.com/unarya/unarya/pkg/parser"
)

// ===============================================
// gRPC Server Entry
// ===============================================
func main() {
	port := os.Getenv("PARSER_PORT")
	if port == "" {
		port = "50053"
	}
	defer os.Exit(1)

	pkg.InitGRPCServer()
	pkg.MinIOConnect()
	log.Printf("[SYSTEM] Parser services running on port %s", port)

	if err := pkg.StartGRPCServer(utils.ToInt(port)); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
