package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

// StartServer ...
func (c *Command) StartServer(ctx context.Context, dir string) error {
	cmdExecDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to get cmd exec dir: %w", err)
	}
	projectFullDir := filepath.Join(cmdExecDir, dir)
	logger.Infof(ctx, "Current project directory: %s", projectFullDir)

	logger.Infof(ctx, "Starting ISO server...")
	if err := c.docker.StartServer(projectFullDir); err != nil {
		return fmt.Errorf("fail start server in docker: %w", err)
	}
	logger.Infof(ctx, "ISO server started")

	return nil
}
