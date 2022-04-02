package command

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

func (c *Command) RulesSync(ctx context.Context, dir string) error {
	svcConfigs, err := c.isoSrv.GetServiceConfigs(ctx)
	if err != nil {
		return fmt.Errorf("fail to get service configs from isoserver: %w", err)
	}

	if err := c.writeRules(ctx, svcConfigs, dir); err != nil {
		return fmt.Errorf("fail to write service configs to directory: %s, err: %w", dir, err)
	}

	return nil
}

func (c *Command) writeRules(ctx context.Context, svcConfigs []models.ServiceConfigDesc, path string) error {
	currDir := filepath.Join(path, config.RulesDir)
	if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
		if os.IsExist(err) {
			logger.Warnf(ctx, "directory %s already exists, skip generating", currDir)
			return nil
		}
		return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
	}

	for _, svcCfg := range svcConfigs {
		if svcCfg.Name == "" {
			return fmt.Errorf("empty service config name for: %s", svcCfg.Host)
		}
		currDir := filepath.Join(currDir, svcCfg.Name)
		if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
			return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
		}
		if err := encodeDataToYamlFile(svcCfg, currDir, config.ServiceConfigFileName); err != nil {
			return fmt.Errorf("fail to encode data to yaml: %w", err)
		}
		if len(svcCfg.GRPCHandlers) == 0 {
			continue
		}
		handlerDir := filepath.Join(currDir, config.ProtoHandlerDirName)
		if err := os.Mkdir(handlerDir, fs.ModePerm); err != nil {
			return fmt.Errorf("fail to create handler dir: %s, err: %w", handlerDir, err)
		}
		for _, handlerCfg := range svcCfg.GRPCHandlers {
			fileName := fmt.Sprintf("%s_%s.yaml", handlerCfg.ServiceName, handlerCfg.MethodName)
			if err := encodeDataToYamlFile(handlerCfg, handlerDir, fileName); err != nil {
				return fmt.Errorf("fail to encode data to yaml: %w", err)
			}
		}
	}

	return nil
}

func encodeDataToYamlFile(data interface{}, currDir, fileName string) error {
	f, err := os.Create(filepath.Join(currDir, fileName))
	if err != nil {
		return fmt.Errorf("fail to create %s: %w", fileName, err)
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("fail to encode data: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("fail to close yaml file: %w", err)
	}
	return nil
}
