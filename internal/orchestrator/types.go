package orchestrator

import "time"

// Request defines a full orchestration request
type Request struct {
	RepositoryURL string
	Branch        string
	Token         string
	SourceType    string // "git", "archive", "url"
}

// ParsedData represents output from the Parser service
type ParsedData struct {
	Language     string
	Dependencies []string
	Metrics      map[string]float64
	Structure    interface{}
}

// AIResult represents the output of a Python AI microservice
type AIResult struct {
	Predictions map[string]float64
	Insights    map[string]string
	ModelUsed   string
}

// Result holds the overall orchestration result
type Result struct {
	CollectorStatus string
	ParserStatus    string
	SecurityStatus  string
	AIStatus        string
	FinalResult     FinalResult
	Duration        time.Duration
}

// FinalResult is the final aggregated outcome
type FinalResult struct {
	Summary     string
	RiskScore   float64
	Insights    map[string]string
	Errors      []string
	CompletedAt time.Time
}

// State represents current pipeline state
type State struct {
	Stage     string
	Status    string // "pending", "running", "success", "failed"
	Timestamp time.Time
}
