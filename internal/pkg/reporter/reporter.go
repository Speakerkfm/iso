package reporter

import (
	"context"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Service interface {
	GetReport(ctx context.Context) (*models.Report, error)
}

type EventService interface {
	GetEvents(ctx context.Context) ([]*models.Event, error)
}

type reporter struct {
	eventService EventService
}

func New(eventService EventService) Service {
	return &reporter{
		eventService: eventService,
	}
}

func (r *reporter) GetReport(ctx context.Context) (*models.Report, error) {
	events, err := r.eventService.GetEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get all events: %w", err)
	}

	report := &models.Report{
		Service: make(map[string]*models.ServiceReport),
	}

	for _, evt := range events {
		if _, ok := report.Service[evt.ServiceName]; !ok {
			report.Service[evt.ServiceName] = &models.ServiceReport{
				Method: make(map[string]*models.MethodReport),
			}
		}

		if _, ok := report.Service[evt.ServiceName].Method[evt.MethodName]; !ok {
			report.Service[evt.ServiceName].Method[evt.MethodName] = &models.MethodReport{
				Stat: &models.MethodStat{},
			}
		}

		methodReport := report.Service[evt.ServiceName].Method[evt.MethodName]

		if evt.IsSuccess {
			methodReport.Stat.SuccessCount++
		} else {
			methodReport.Stat.ErrorCount++
		}
	}

	return report, nil
}
