package fetcher

import (
	"context"
	"strings"

	"github.com/Speakerkfm/iso/internal/pkg/fetcher/local"
	"github.com/Speakerkfm/iso/internal/pkg/fetcher/web"
)

type FileFetcher interface {
	FetchFile(ctx context.Context, filePath string) ([]byte, error)
}

type Fetcher struct {
	webFetcher   FileFetcher
	localFetcher FileFetcher
}

func New() *Fetcher {
	webFetcher := web.New()
	localFetcher := local.New()

	return &Fetcher{
		webFetcher:   webFetcher,
		localFetcher: localFetcher,
	}
}

// FetchFile загружает файл по указанному пути
func (f *Fetcher) FetchFile(ctx context.Context, filePath string) ([]byte, error) {
	if strings.HasPrefix(filePath, "http") {
		return f.webFetcher.FetchFile(ctx, filePath)
	}
	return f.localFetcher.FetchFile(ctx, filePath)
}
