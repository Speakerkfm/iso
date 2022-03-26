package command

import (
	"context"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

// Root - дефолтная команда
func (c *Command) Root(ctx context.Context) error {
	logger.Info(ctx, "It works")
	return nil
}
