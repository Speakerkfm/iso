package batcher

import (
	"context"
	"time"

	"go.uber.org/atomic"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type batch struct {
	ctx     context.Context
	number  int
	events  []*models.Event
	eventCh <-chan *models.Event
	closeCh <-chan struct{}

	eventRepo eventRepository

	flushInterval     *atomic.Duration
	flushEventsAmount *atomic.Int64
}

func newBatch(
	ctx context.Context,
	number int,
	eventCh <-chan *models.Event,
	closeCh <-chan struct{},
	eventRepo eventRepository,
	flushInterval *atomic.Duration,
	flushEventsAmount *atomic.Int64) *batch {
	return &batch{
		ctx:               ctx,
		number:            number,
		events:            make([]*models.Event, 0, flushEventsAmount.Load()),
		eventCh:           eventCh,
		closeCh:           closeCh,
		eventRepo:         eventRepo,
		flushInterval:     flushInterval,
		flushEventsAmount: flushEventsAmount,
	}
}

func (b *batch) observeFlush() {
	logger.Infof(b.ctx, "batcher %d started observing flush", b.number)

	t := time.NewTimer(b.flushInterval.Load())
	for {
		select {
		case <-b.closeCh:
			logger.Infof(b.ctx, "batcher %d got terminated signal", b.number)
			b.flush()
			return
		case <-t.C:
			b.flush()
			t.Reset(b.flushInterval.Load())
		case e, ok := <-b.eventCh:
			if !ok {
				// channel is closed - flush and stop observing
				b.flush()
				return
			}

			b.events = append(b.events, e)
			if len(b.events) >= int(b.flushEventsAmount.Load()) {
				b.flush()
				t.Reset(b.flushInterval.Load())
			}
		}
	}
}

func (b *batch) flush() {
	if len(b.events) == 0 {
		return
	}
	if err := b.eventRepo.SaveBatch(b.ctx, b.events); err != nil {
		logger.Errorf(b.ctx, "error flushing %d batch", b.number, "size", len(b.events), "err", err)
	}

	b.events = make([]*models.Event, 0, b.flushEventsAmount.Load())
}
