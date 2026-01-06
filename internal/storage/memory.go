package storage

import (
	"sync"

	"sovana/internal/metrics"
)

// MemoryStore keeps metrics in memory in a thread-safe manner.
type MemoryStore struct {
	mu      sync.RWMutex
	metrics []metrics.Metric
}

// NewMemoryStore initializes an empty store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// Add appends a new metric snapshot.
func (s *MemoryStore) Add(m metrics.Metric) {
	s.mu.Lock()
	s.metrics = append(s.metrics, m)
	s.mu.Unlock()
}

// All returns a copy of all stored metrics.
func (s *MemoryStore) All() []metrics.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]metrics.Metric, len(s.metrics))
	n := copy(items, s.metrics)
	return items[:n]
}
