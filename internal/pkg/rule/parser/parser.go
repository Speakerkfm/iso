package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

const (
	protoHandlerDirName   = "grpc"
	serviceConfigFileName = "service.yaml"
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
				currDir := path.Join(currDir, protoHandlerDirName)
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

// GenerateRules генерирует правила из конфигов сервисов
func (p *parser) GenerateRules(svcConfigs []models.ServiceConfigDesc) []*models.Rule {
	var res []*models.Rule
	for _, svcConfig := range svcConfigs {
		for _, handlerCfg := range svcConfig.GRPCHandlers {
			for _, handlerRule := range handlerCfg.Rules {
				r := &models.Rule{
					Conditions: []models.Condition{
						{
							Key:   models.FieldHost,
							Value: svcConfig.Host,
						},
						{
							Key:   models.FieldServiceName,
							Value: handlerCfg.ServiceName,
						},
						{
							Key:   models.FieldMethodName,
							Value: handlerCfg.MethodName,
						},
					},
					HandlerConfig: &models.HandlerConfig{
						ResponseDelay: handlerRule.Response.Delay,
						MessageData:   json.RawMessage(handlerRule.Response.Data),
					},
				}
				for _, cond := range handlerRule.Conditions {
					r.Conditions = append(r.Conditions, models.Condition{
						Key:   cond.Key,
						Value: cond.Value,
					})
				}
				res = append(res, r)
			}
		}
	}
	return res
}

func isGRPCRulesDir(file fs.FileInfo) bool {
	return file.IsDir() && file.Name() == protoHandlerDirName
}

func isServiceConfig(file fs.FileInfo) bool {
	return !file.IsDir() && file.Name() == serviceConfigFileName
}
