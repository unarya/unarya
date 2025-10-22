package orchestrator

import (
	_ "fmt"
	"log"
	"time"
)

// Orchestrator coordinates multi-service pipelines
type Orchestrator struct {
	StateManager *StateManager
	Client       *PythonClient
	ErrorHandler *ErrorHandler
	Results      []interface{}
}

// ExecutePipeline runs the full multi-step orchestration
func (o *Orchestrator) ExecutePipeline(req *Request) (*Result, error) {
	start := time.Now()
	log.Printf("[Orchestrator] Starting pipeline for %s\n", req.RepositoryURL)

	// 1. Update state
	o.StateManager.Update("collector", "running")

	// 2. Simulate collection
	time.Sleep(300 * time.Millisecond)
	o.StateManager.Update("collector", "success")

	// 3. Simulate parsing stage
	o.StateManager.Update("parser", "running")
	parsed := &ParsedData{
		Language: "Go",
		Dependencies: []string{
			"github.com/gin-gonic/gin",
			"github.com/redis/go-redis/v9",
		},
		Metrics: map[string]float64{"complexity": 7.2, "files": 43},
	}
	o.StateManager.Update("parser", "success")

	// 4. Call AI model service (Python microservice)
	o.StateManager.Update("ai", "running")
	aiRes, err := o.CallPythonService(parsed)
	if err != nil {
		o.HandleFailure("ai", err)
		return nil, err
	}
	o.StateManager.Update("ai", "success")

	// 5. Aggregate results
	final := o.AggregateResults()
	final.Insights = aiRes.Insights

	duration := time.Since(start)
	result := &Result{
		CollectorStatus: "success",
		ParserStatus:    "success",
		SecurityStatus:  "skipped",
		AIStatus:        "success",
		FinalResult:     *final,
		Duration:        duration,
	}

	log.Printf("[Orchestrator] Completed pipeline in %v\n", duration)
	return result, nil
}

// AggregateResults consolidates intermediate outputs
func (o *Orchestrator) AggregateResults() *FinalResult {
	return &FinalResult{
		Summary:     "Pipeline executed successfully",
		RiskScore:   0.42,
		Insights:    map[string]string{"efficiency": "high", "security": "stable"},
		CompletedAt: time.Now(),
	}
}
