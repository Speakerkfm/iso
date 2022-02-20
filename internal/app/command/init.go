package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

const (
	defaultPath = "."

	configFileName = "spec.yaml"
)

func (c *Command) Init(ctx context.Context, path string) error {
	specData, err := c.gen.GenerateSpecificationData()
	if err != nil {
		return fmt.Errorf("fail to generate config data: %w", err)
	}

	if path == "" {
		path = defaultPath
	}

	if err := ioutil.WriteFile(filepath.Join(path, configFileName), specData, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to save spec data to file: %w", err)
	}

	logger.Info(ctx, "Project initialized")

	return nil
}
