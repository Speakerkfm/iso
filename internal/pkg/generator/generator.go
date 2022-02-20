package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/bxcodec/faker/v3"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	public_models "github.com/Speakerkfm/iso/pkg/models"
)

// Generator генерирует файлы по шаблонам
type Generator interface {
	GenerateSpecificationData() ([]byte, error)
	GeneratePluginData(pluginSpec models.PluginDesc) ([]byte, error)
	GenerateServiceConfigs(spec models.ServiceSpecification, svcProvider public_models.ServiceProvider) ([]models.ServiceConfigDesc, error)
}

type generator struct {
}

// New создает объект генератора
func New() Generator {
	return &generator{}
}

// GenerateConfigData генерирует пример файла конфигурации
func (g *generator) GenerateSpecificationData() ([]byte, error) {
	return configTemplateExample, nil
}

// GeneratePluginData генерирует данные для .go файла плагина, который возвращает описание структур
func (g *generator) GeneratePluginData(pluginSpec models.PluginDesc) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	if err := implTemplate.Execute(buff, pluginSpec); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (g *generator) GenerateServiceConfigs(spec models.ServiceSpecification, svcProvider public_models.ServiceProvider) ([]models.ServiceConfigDesc, error) {
	dataPathProtoService := make(map[string]*public_models.ProtoService)

	for _, svc := range svcProvider.GetList() {
		protoFileName := filepath.Base(svc.ProtoPath)
		dataPathProtoService[protoFileName] = svc
	}

	var svcExamples []models.ServiceConfigDesc

	for _, dep := range spec.ExternalDependencies {
		svcDesc := models.ServiceConfigDesc{
			Host: dep.Host,
			Name: dep.Name,
		}
		for _, protoDepPath := range dep.ProtoPaths {
			protoFileName := filepath.Base(protoDepPath)
			protoSvc, ok := dataPathProtoService[protoFileName]
			if !ok {
				logger.Errorf(context.Background(), "fail to get proto svc by proto file: %s", protoFileName)
				continue
			}

			for _, protoHandler := range protoSvc.Methods {
				handlerDesc, err := generateProtoHandlerConfig(protoSvc, protoHandler)
				if err != nil {
					logger.Errorf(context.Background(), "fail to generate proto handler cfg: %s", err.Error())
					continue
				}
				svcDesc.GRPCHandlers = append(svcDesc.GRPCHandlers, handlerDesc)
			}
		}
		svcExamples = append(svcExamples, svcDesc)
	}

	return svcExamples, nil
}

func generateProtoHandlerConfig(protoSvc *public_models.ProtoService, protoHandler public_models.ProtoMethod) (models.HandlerConfigDesc, error) {
	respStruct := protoHandler.ResponseStruct
	if err := faker.FakeData(respStruct); err != nil {
		return models.HandlerConfigDesc{}, fmt.Errorf("fail to generate fake data svc: %s, handler: %s, err: %w", protoSvc.Name, protoHandler.Name, err)
	}
	respData, err := json.MarshalIndent(respStruct, "", "	")
	if err != nil {
		return models.HandlerConfigDesc{}, fmt.Errorf("fail to marshal resp example svc: %s, handler: %s, err: %w", protoSvc.Name, protoHandler.Name, err)
	}
	return models.HandlerConfigDesc{
		ServiceName: protoSvc.Name,
		MethodName:  protoHandler.Name,
		Rules: []models.RuleDesc{
			{
				Conditions: []models.HandlerConditionDesc{
					{
						Key:   "header.x-request-id",
						Value: "*",
					},
				},
				Response: models.HandlerResponseDesc{
					Data:  string(respData),
					Delay: 5 * time.Millisecond,
				},
			},
		},
	}, nil
}
