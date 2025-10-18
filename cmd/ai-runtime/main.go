package main

import (
	"fmt"
	"log"
	"net"

	"github.com/unarya/unarya/lib/proto/pb/aipb"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	runtimeSrv := &RuntimeServer{}

	aipb.RegisterAIInferenceServer(server, runtimeSrv)

	// Load default model on startup
	if err := runtimeSrv.LoadModel("/app/models/model.onnx"); err != nil {
		log.Fatalf("❌ Failed to load model: %v", err)
	}

	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("[AI Runtime] 🚀 Serving on port 6000 (GPU mode)")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
