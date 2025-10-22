package orchestrator

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

// PythonClient handles communication with Python-based AI microservice
type PythonClient struct {
	conn   *grpc.ClientConn
	target string
}

// NewPythonClient initializes the gRPC client
func NewPythonClient(target string) (*PythonClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Python service: %w", err)
	}
	return &PythonClient{conn: conn, target: target}, nil
}

// CallPythonService executes the AI inference request
func (o *Orchestrator) CallPythonService(data *ParsedData) (*AIResult, error) {
	log.Printf("[Orchestrator] Calling Python service for language=%s\n", data.Language)

	// In reality you'd call generated proto client, e.g. aipb.AIServiceClient
	time.Sleep(200 * time.Millisecond)

	return &AIResult{
		Predictions: map[string]float64{"quality_score": 0.91},
		Insights:    map[string]string{"refactor": "low priority"},
		ModelUsed:   "ai-model-v2",
	}, nil
}
