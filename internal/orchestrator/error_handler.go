package orchestrator

import (
	"fmt"
	"log"
	"time"
)

// ErrorHandler manages retries and error logging
type ErrorHandler struct {
	MaxRetries int
	Delay      time.Duration
}

// NewErrorHandler initializes a retry handler
func NewErrorHandler(maxRetries int, delay time.Duration) *ErrorHandler {
	return &ErrorHandler{
		MaxRetries: maxRetries,
		Delay:      delay,
	}
}

// HandleFailure handles error recovery and retry logic
func (o *Orchestrator) HandleFailure(stage string, err error) error {
	log.Printf("[ErrorHandler] Failure at stage %s: %v\n", stage, err)
	for i := 1; i <= o.ErrorHandler.MaxRetries; i++ {
		log.Printf("[ErrorHandler] Retry %d/%d for stage %s\n", i, o.ErrorHandler.MaxRetries, stage)
		time.Sleep(o.ErrorHandler.Delay)

		if i == o.ErrorHandler.MaxRetries {
			return fmt.Errorf("stage %s failed after %d retries: %w", stage, o.ErrorHandler.MaxRetries, err)
		}
	}
	return nil
}
