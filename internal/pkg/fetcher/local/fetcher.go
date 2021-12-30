package local

import (
	"context"
	"io/ioutil"
)

type Fetcher struct {
}

func New() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) FetchFile(ctx context.Context, filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
