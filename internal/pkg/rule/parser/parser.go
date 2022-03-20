package parser

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type parser struct {
}

func New() *parser {
	return &parser{}
}

// ParseDirectory парсит конфигурации сервисов из указанной директории во внутренние структуры
func (p *parser) ParseDirectory(ctx context.Context, directoryPath string) ([]models.ServiceConfigDesc, error) {
	var res []models.ServiceConfigDesc

	serviceFiles, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("fail to read dir: %s: %w", directoryPath, err)
	}

	for _, serviceFile := range serviceFiles {
		if !serviceFile.IsDir() {
			return nil, fmt.Errorf("invalid service dir")
		}

		currDir := path.Join(directoryPath, serviceFile.Name())
		svc := models.ServiceConfigDesc{}

		serviceConfigFiles, err := ioutil.ReadDir(currDir)
		if err != nil {
			return nil, fmt.Errorf("fail to read dir: %s: %w", currDir, err)
		}

		for _, serviceConfigFile := range serviceConfigFiles {
			if isGRPCRulesDir(serviceConfigFile) {
				currDir := path.Join(currDir, config.ProtoHandlerDirName)
				handlerConfigFiles, err := ioutil.ReadDir(currDir)
				if err != nil {
					return nil, fmt.Errorf("fail to read dir: %s: %w", currDir, err)
				}

				for _, handlerConfigFile := range handlerConfigFiles {
					handlerConfigData, err := os.Open(path.Join(currDir, handlerConfigFile.Name()))
					if err != nil {
						return nil, fmt.Errorf("fail to read file: %s from dir: %s: %w", handlerConfigFile.Name(), currDir, err)
					}

					handlerConfig := models.HandlerConfigDesc{}

					if err := yaml.NewDecoder(handlerConfigData).Decode(&handlerConfig); err != nil {
						return nil, fmt.Errorf("fail to decode handler config: %s from dir: %s: %w", handlerConfigFile.Name(), currDir, err)
					}
					handlerConfigData.Close()

					svc.GRPCHandlers = append(svc.GRPCHandlers, handlerConfig)
				}
			}

			if isServiceConfig(serviceConfigFile) {
				svcFile, err := os.Open(path.Join(currDir, serviceConfigFile.Name()))
				if err != nil {
					return nil, fmt.Errorf("fail to read file: %s from dir: %s: %w", serviceConfigFile.Name(), currDir, err)
				}
				if err := yaml.NewDecoder(svcFile).Decode(&svc); err != nil {
					return nil, fmt.Errorf("fail to decode service config: %s from dir: %s: %w", serviceConfigFile.Name(), currDir, err)
				}
				svcFile.Close()
			}
		}

		res = append(res, svc)
	}

	return res, nil
}

func isGRPCRulesDir(file fs.FileInfo) bool {
	return file.IsDir() && file.Name() == config.ProtoHandlerDirName
}

func isServiceConfig(file fs.FileInfo) bool {
	return !file.IsDir() && file.Name() == config.ServiceConfigFileName
}
