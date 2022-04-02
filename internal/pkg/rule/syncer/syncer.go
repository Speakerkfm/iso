package syncer

import (
	"context"
	"sync"
	"time"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Source interface {
	GetRules(ctx context.Context) ([]*models.Rule, error)
}

type Target interface {
	UpdateRuleTree(rules []*models.Rule)
}

type Syncer struct {
	ctx      context.Context
	src      Source
	target   Target
	interval time.Duration
	wg       *sync.WaitGroup
	done     chan struct{}
}

func New(ctx context.Context, src Source, target Target, interval time.Duration) *Syncer {
	return &Syncer{
		ctx:      ctx,
		src:      src,
		target:   target,
		interval: interval,
		wg:       &sync.WaitGroup{},
		done:     make(chan struct{}),
	}
}

// Start ...
func (s *Syncer) Start() {
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		logger.Infof(s.ctx, "syncer started")
		t := time.NewTicker(s.interval)
		for {
			select {
			case <-s.ctx.Done():
				logger.Infof(s.ctx, "syncer stopped")
				return
			case <-t.C:
				s.sync()
				logger.Infof(s.ctx, "rules synced")
			case <-s.done:
				logger.Infof(s.ctx, "syncer stopped")
				return
			}
		}
	}()
}

func (s *Syncer) Stop() {
	close(s.done)
	s.wg.Wait()
}

func (s *Syncer) sync() {
	rules, err := s.src.GetRules(s.ctx)
	if err != nil {
		logger.Errorf(s.ctx, "fail to get rules from source: %s", err.Error())
		return
	}

	s.target.UpdateRuleTree(rules)
}
