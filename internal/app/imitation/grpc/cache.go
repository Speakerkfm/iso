package grpc

import (
	"context"
	"sync"
)

type responseStore struct {
	mu   sync.RWMutex
	data map[string]*Response
}

func newResponseStore() *responseStore {
	return &responseStore{
		data: make(map[string]*Response),
	}
}

func (st *responseStore) Get(ctx context.Context, key string) (*Response, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	resp, ok := st.data[key]

	return resp, ok
}

func (st *responseStore) Set(ctx context.Context, key string, resp *Response) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.data[key] = resp
}
