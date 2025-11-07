package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/unarya/unarya/lib/proto/pb/aipb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:6000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("❌ Cannot connect to ML service: %v", err)
	}
	defer conn.Close()

	client := aipb.NewAIServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &aipb.AIAnalyzeRequest{
		Language:      "javascript",
		CodeStructure: "function hello(){ console.log('hi') }",
	}

	res, err := client.AnalyzeCode(ctx, req)
	if err != nil {
		log.Fatalf("❌ ML model error: %v", err)
	}

	fmt.Println("✅ ML model response:")
	fmt.Printf("  → insights: %s\n", res.GetInsights())
	fmt.Printf("  → confidence: %s\n", res.GetConfidence())
}
