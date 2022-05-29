package events

import (
	"context"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type batcher interface {
	Append(ctx context.Context, e *models.Event) error
}

type eventRepository interface {
	GetEvents(ctx context.Context) ([]*models.Event, error)
}

type Service struct {
	batcher   batcher
	eventRepo eventRepository
}

func New(batcher batcher, eventRepo eventRepository) *Service {
	return &Service{
		batcher:   batcher,
		eventRepo: eventRepo,
	}
}

func (svc *Service) PushEvent(ctx context.Context, serviceName, methodName string, isSuccess bool) error {
	return svc.batcher.Append(ctx, &models.Event{
		ServiceName: serviceName,
		MethodName:  methodName,
		IsSuccess:   isSuccess,
	})
}

func (svc *Service) GetEvents(ctx context.Context) ([]*models.Event, error) {
	return svc.eventRepo.GetEvents(ctx)
}
