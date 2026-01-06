package storage

import (
	"context"
	"sort"
	"sync"
	"time"

	"sovana/internal/model"
)

// MemoryStore provides an in-memory thread-safe store.
type MemoryStore struct {
	mu     sync.RWMutex
	items  map[int64]model.Item
	nextID int64
}

// NewMemoryStore initializes an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items:  make(map[int64]model.Item),
		nextID: 1,
	}
}

// List returns all items sorted by ID.
func (s *MemoryStore) List(ctx context.Context) ([]model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(s.items))
	for id := range s.items {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	result := make([]model.Item, 0, len(s.items))
	for _, id := range ids {
		result = append(result, s.items[id])
	}

	return result, nil
}

// Create stores a new item and returns it with an assigned ID.
func (s *MemoryStore) Create(ctx context.Context, item model.Item) (model.Item, error) {
	if err := ctx.Err(); err != nil {
		return model.Item{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item.ID = s.nextID
	item.CreatedAt = time.Now().UTC()

	s.items[item.ID] = item
	s.nextID++

	return item, nil
}

// Delete removes an item by ID.
func (s *MemoryStore) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}

	delete(s.items, id)
	return nil
}
