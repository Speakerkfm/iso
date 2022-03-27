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

func (c *Command) RulesDownload(ctx context.Context, dir string) error {
	return nil
}

func (c *Command) generateRuleExamples(ctx context.Context, spec models.ServiceSpecification, path string) error {
	svcProvider, err := loadPluginData(ctx, filepath.Join(path, config.PluginFileName))
	if err != nil {
		return fmt.Errorf("fail to load plugin data: %w", err)
	}

	svcConfigs, err := c.gen.GenerateServiceConfigs(spec, svcProvider)
	if err != nil {
		return fmt.Errorf("fail to generate svc cfgs: %w", err)
	}

	currDir := filepath.Join(path, config.RulesDir)
	if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
		if os.IsExist(err) {
			logger.Warnf(ctx, "directory %s already exists, skip generating", currDir)
			return nil
		}
		return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
	}

	for _, svcCfg := range svcConfigs {
		currDir := filepath.Join(currDir, svcCfg.Name)
		if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
			return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
		}
		if err := encodeDataToYamlFile(svcCfg, currDir, config.ServiceConfigFileName); err != nil {
			return fmt.Errorf("fail to encode data to yaml: %w", err)
		}
		if len(svcCfg.GRPCHandlers) > 0 {
			currDir := filepath.Join(currDir, config.ProtoHandlerDirName)
			if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
				return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
			}
			for _, handlerCfg := range svcCfg.GRPCHandlers {
				fileName := fmt.Sprintf("%s_%s.yaml", handlerCfg.ServiceName, handlerCfg.MethodName)
				if err := encodeDataToYamlFile(handlerCfg, currDir, fileName); err != nil {
					return fmt.Errorf("fail to encode data to yaml: %w", err)
				}
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
