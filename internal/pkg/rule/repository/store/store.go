package store

import (
	"context"
	"sync"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Store struct {
	mu         sync.RWMutex
	svcConfigs []models.ServiceConfigDesc
}

func New() *Store {
	return &Store{}
}

func (st *Store) GetServiceConfigs(ctx context.Context) ([]models.ServiceConfigDesc, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.svcConfigs, nil
}

func (st *Store) SaveServiceConfigs(ctx context.Context, svcConfigs []models.ServiceConfigDesc) error {
	st.mu.Lock()
	st.svcConfigs = svcConfigs
	st.mu.Unlock()

	return nil
}
