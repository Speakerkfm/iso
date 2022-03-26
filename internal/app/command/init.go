package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

// Init - создает проект в указанной директории
func (c *Command) Init(ctx context.Context, dir string) error {
	specData, err := c.gen.GenerateSpecificationData()
	if err != nil {
		return fmt.Errorf("fail to generate config data: %w", err)
	}

	if dir == "" {
		dir = config.DefaultProjectDir
	}

	if err := ioutil.WriteFile(filepath.Join(dir, config.SpecificationFileName), specData, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to save spec data to file: %w", err)
	}

	logger.Info(ctx, "Project initialized")

	return nil
}
