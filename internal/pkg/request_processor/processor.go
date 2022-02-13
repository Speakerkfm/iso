package request_processor

import (
	"context"
	"fmt"
	"time"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Processor процессит все пришедшие запросы в имитирующий сервис
type Processor interface {
	Process(ctx context.Context, req *models.Request) (*models.Response, error)
}

type RuleManager interface {
	GetHandlerConfig(ctx context.Context, req *models.Request) (*models.HandlerConfig, error)
}

type processor struct {
	ruleManager RuleManager
}

// New создает новый процессор
func New(ruleManager RuleManager) Processor {
	return &processor{
		ruleManager: ruleManager,
	}
}

// Process обрабатывает пришедший запрос
func (p *processor) Process(ctx context.Context, req *models.Request) (*models.Response, error) {
	cfg, err := p.ruleManager.GetHandlerConfig(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to get rule for request: %w", err)
	}

	if err := waitDelay(ctx, cfg.ResponseDelay); err != nil {
		return nil, fmt.Errorf("fail to wait delay: %w", err)
	}

	return &models.Response{
		Message: cfg.MessageData,
	}, nil
}

func waitDelay(ctx context.Context, delay time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
	}
	return nil
}

// `{"exists":true}`
// `{"user":{"id":10,"name":"kek"}}`
