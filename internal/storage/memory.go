package storage

import (
	"sync"

	"sovana/internal/model"
)

// MemoryStorage keeps transactions and subscriptions in-memory safely.
type MemoryStorage struct {
	mu            sync.RWMutex
	transactions  []model.Transaction
	subscriptions []model.Subscription
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// Save overwrites the current dataset.
func (s *MemoryStorage) Save(transactions []model.Transaction, subscriptions []model.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions = append([]model.Transaction(nil), transactions...)
	s.subscriptions = append([]model.Subscription(nil), subscriptions...)
}

// Subscriptions returns a copy of the detected subscriptions.
func (s *MemoryStorage) Subscriptions() []model.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Subscription(nil), s.subscriptions...)
}
