package web

import (
	"context"
)

type Fetcher struct {
}

func New() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) FetchFile(ctx context.Context, filePath string) ([]byte, error) {
	return nil, nil
}
