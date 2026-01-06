package storage

import (
	"context"
	"errors"

	"sovana/internal/model"
)

// Store defines the persistence contract for Items.
type Store interface {
	List(ctx context.Context) ([]model.Item, error)
	Create(ctx context.Context, item model.Item) (model.Item, error)
	Delete(ctx context.Context, id int64) error
}

// ErrNotFound is returned when a requested item does not exist.
var ErrNotFound = errors.New("item not found")
