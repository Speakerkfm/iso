package request_processor

import (
	"context"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Processor процессит все пришедшие запросы в имитирующий сервис
type Processor struct {
}

// New создает новый процессор
func New() *Processor {
	return &Processor{}
}

// Process обрабатывает пришедший запрос
func (p *Processor) Process(ctx context.Context, req *models.Request) (*models.Response, error) {
	return nil, nil
}
