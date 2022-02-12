package request_processor

import (
	"context"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Processor процессит все пришедшие запросы в имитирующий сервис
type Processor interface {
	Process(ctx context.Context, req *models.Request) (*models.Response, error)
}

type RuleManager interface {
	GetRule(ctx context.Context, req *models.Request) (*models.Rule, error)
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
	rule, err := p.ruleManager.GetRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to get rule for request: %w", err)
	}

	return &models.Response{
		Message: rule.MessageData,
	}, nil
}

// `{"exists":true}`
// `{"user":{"id":10,"name":"kek"}}`
