package batcher

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/atomic"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type eventRepository interface {
	SaveBatch(ctx context.Context, events []*models.Event) error
}

// Batcher for sending events in batch to clickhouse
type Batcher struct {
	ctx     context.Context
	wg      *sync.WaitGroup
	enabled *atomic.Bool
	eventCh chan<- *models.Event
	closeCh chan struct{}
	batches []*batch
}

// New creates new batcher
func New(ctx context.Context,
	enabled *atomic.Bool,
	eventRepo eventRepository,
	batchCount int,
	flushInterval *atomic.Duration,
	flushEventsAmount *atomic.Int64,
	eventBuffSize int,
) *Batcher {
	wg := &sync.WaitGroup{}
	eventCh := make(chan *models.Event, eventBuffSize)
	closeCh := make(chan struct{})

	batcher := &Batcher{
		ctx:     ctx,
		wg:      wg,
		enabled: enabled,
		eventCh: eventCh,
		closeCh: closeCh,
		batches: make([]*batch, 0, batchCount),
	}

	for num := 0; num < batchCount; num++ {
		b := newBatch(ctx,
			num,
			eventCh,
			closeCh,
			eventRepo,
			flushInterval,
			flushEventsAmount)
		batcher.batches = append(batcher.batches, b)

		wg.Add(1)
		go func() {
			defer wg.Done()
			b.observeFlush()
		}()
	}

	return batcher
}

// Close closes batcher
func (b *Batcher) Close() error {
	close(b.closeCh)
	b.wg.Wait()

	logger.Info(b.ctx, "batcher stopped")

	return nil
}

// Append adds event to batch
func (b *Batcher) Append(ctx context.Context, e *models.Event) error {
	if !b.enabled.Load() {
		return nil
	}

	select {
	case <-b.closeCh:
		return fmt.Errorf("batcher is closed")
	case <-ctx.Done():
		logger.Warnf(ctx, "batcher queue is full - you need to think about another configuration")
		return ctx.Err()
	case b.eventCh <- e:
	}
	return nil
}
