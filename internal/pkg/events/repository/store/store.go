package store

import (
	"context"
	"sync"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Store struct {
	mu     sync.RWMutex
	events []*models.Event
}

func New() *Store {
	return &Store{}
}

func (st *Store) GetEvents(ctx context.Context) ([]*models.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.events, nil
}

func (st *Store) SaveBatch(ctx context.Context, events []*models.Event) error {
	st.mu.Lock()
	st.events = append(st.events, events...)
	st.mu.Unlock()

	return nil
}
