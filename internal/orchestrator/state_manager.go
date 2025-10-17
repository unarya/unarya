package orchestrator

import (
	"log"
	"sync"
	"time"
)

// StateManager tracks the status of each pipeline stage
type StateManager struct {
	mu     sync.Mutex
	states map[string]State
}

func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[string]State),
	}
}

// Update sets the state for a given stage
func (s *StateManager) Update(stage, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := State{
		Stage:     stage,
		Status:    status,
		Timestamp: time.Now(),
	}
	s.states[stage] = state
	log.Printf("[StateManager] Stage=%s, Status=%s\n", stage, status)
}

// Get retrieves state for a specific stage
func (s *StateManager) Get(stage string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[stage]
	return st, ok
}

// Snapshot returns a copy of the entire state map
func (s *StateManager) Snapshot() map[string]State {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make(map[string]State)
	for k, v := range s.states {
		copy[k] = v
	}
	return copy
}
