package command

import (
	"context"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

func (c *Command) Root(ctx context.Context) error {
	logger.Info(ctx, "It works")
	return nil
}
