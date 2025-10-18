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
	// Kết nối tới server Python (không mã hóa TLS, dùng insecure)
	conn, err := grpc.NewClient("localhost:6000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("❌ Không thể kết nối gRPC server: %v", err)
	}
	defer conn.Close()

	client := aipb.NewAIServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &aipb.AIAnalyzeRequest{
		Language:      "python",
		CodeStructure: "def hello(): print('hi')",
	}

	res, err := client.AnalyzeCode(ctx, req)
	if err != nil {
		log.Fatalf("❌ RPC lỗi: %v", err)
	}

	fmt.Println("✅ gRPC response:")
	fmt.Printf("  → insights: %s\n", res.GetInsights())
	fmt.Printf("  → confidence: %s\n", res.GetConfidence())
}
