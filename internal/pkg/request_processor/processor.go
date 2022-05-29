package request_processor

import (
	"context"
	"fmt"
	"time"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/metrics"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Processor процессит все пришедшие запросы в имитирующий сервис
type Processor interface {
	Process(ctx context.Context, req models.Request) (*models.Response, error)
}

// RuleManager выбирает правило, в соответствии с которым нужно обработать пришедший запрос
type RuleManager interface {
	GetRule(ctx context.Context, req models.Request) (*models.Rule, error)
}

// EventService сохраняет события об обработанных запросах для формирования отчетов
type EventService interface {
	PushEvent(ctx context.Context, serviceName, methodName, ruleName string) error
}

type processor struct {
	ruleManager  RuleManager
	eventService EventService
}

// New создает новый процессор запросов
func New(ruleManager RuleManager, eventService EventService) Processor {
	return &processor{
		ruleManager:  ruleManager,
		eventService: eventService,
	}
}

// Process главный обработчик пришедшего запроса
func (p *processor) Process(ctx context.Context, req models.Request) (*models.Response, error) {
	started := time.Now()

	rule, err := p.ruleManager.GetRule(ctx, req)
	defer writeRequestProcessTimeMetric(started, rule)
	if err != nil {
		return nil, fmt.Errorf("fail to get rule for request: %w", err)
	}

	respCfg := rule.HandlerConfig

	if err := waitDelay(ctx, req.GetHandledAt(), respCfg.ResponseDelay); err != nil {
		return nil, fmt.Errorf("fail to wait delay: %w", err)
	}

	if err := p.eventService.PushEvent(ctx, respCfg.ServiceName, respCfg.MethodName, rule.Name); err != nil {
		logger.Errorf(ctx, "fail to push event: %s", err.Error())
	}

	return &models.Response{
		ID:      rule.ID,
		Message: respCfg.MessageData,
		Error:   respCfg.Error,
	}, nil
}

func waitDelay(ctx context.Context, handledAt time.Time, delay time.Duration) error {
	if time.Since(handledAt) >= delay {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
	}
	return nil
}

func writeRequestProcessTimeMetric(started time.Time, rule *models.Rule) {
	if rule == nil {
		return
	}
	metrics.RequestProcessingTimeHistogramVec.WithLabelValues(
		rule.HandlerConfig.ServiceName,
		rule.HandlerConfig.MethodName).
		Observe(time.Since(started).Seconds())
}
