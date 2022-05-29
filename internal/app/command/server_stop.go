package command

import (
	"context"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

// StopServer ...
func (c *Command) StopServer(ctx context.Context, dockerEnabled bool) error {
	if dockerEnabled {
		logger.Infof(ctx, "Stopping ISO server in docker...")
		if err := c.docker.StopServer(); err != nil {
			return fmt.Errorf("fail stop server in docker: %w", err)
		}
		logger.Infof(ctx, "ISO server stopped")

		return nil
	}

	return nil
}
