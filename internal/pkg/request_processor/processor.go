package request_processor

import (
	"context"
	"encoding/json"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Processor процессит все пришедшие запросы в имитирующий сервис
type Processor interface {
	Process(ctx context.Context, req *models.Request) (*models.Response, error)
}

type processor struct {
}

// New создает новый процессор
func New() Processor {
	return &processor{}
}

// Process обрабатывает пришедший запрос
func (p *processor) Process(ctx context.Context, req *models.Request) (*models.Response, error) {
	return &models.Response{
		Message: json.RawMessage(`{"exists":true}`),
	}, nil
}
