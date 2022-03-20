package fetcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type FileFetcher interface {
	FetchFile(ctx context.Context, filePath string) ([]byte, error)
}

type fetcher struct {
}

func New() *fetcher {
	return &fetcher{}
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
	if !strings.HasPrefix(filePath, "http") {
		filePath = fmt.Sprintf("http://%s", filePath)
	}
	resp, err := http.Get(filePath)
	if err != nil {
		return nil, fmt.Errorf("fail to get file from web: %w", err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func isFilePathWeb(path string) bool {
	return !strings.HasPrefix(path, "/")
}
