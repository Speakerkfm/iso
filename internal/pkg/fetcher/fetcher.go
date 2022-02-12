package fetcher

import (
	"context"
	"io/ioutil"
	"strings"
)

type FileFetcher interface {
	FetchFile(ctx context.Context, filePath string) ([]byte, error)
}

type fetcher struct {
}

func New() *fetcher {
	return &fetcher{
	}
}

// FetchFile загружает файл по указанному пути
func (f *fetcher) FetchFile(ctx context.Context, filePath string) ([]byte, error) {
	if isFilePathWeb(filePath) {
		return f.fetchWebFile(ctx, filePath)
	}
	return f.fetchLocalFile(ctx, filePath)
}

func (f *fetcher) fetchLocalFile(ctx context.Context, filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func (f *fetcher) fetchWebFile(ctx context.Context, filePath string) ([]byte, error) {
	return nil, nil
}

func isFilePathWeb(path string) bool {
	return strings.HasPrefix(path, "http")
}
