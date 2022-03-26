package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/generator/adapter/faker"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	public_models "github.com/Speakerkfm/iso/pkg/models"
)

// Generator генерирует файлы по шаблонам
type Generator interface {
	GenerateSpecificationData() ([]byte, error)
	GeneratePluginData(pluginSpec models.PluginDesc) ([]byte, error)
	GenerateServiceConfigs(spec models.ServiceSpecification, svcProvider public_models.ServiceProvider) ([]models.ServiceConfigDesc, error)
	GenerateRules(svcConfigs []models.ServiceConfigDesc) []*models.Rule
	GenerateReverseProxyConfigData() ([]byte, error)
}

// fakeData заполняет сущность рандомными данными
type fakeDataFunc func(a interface{}) error

type generator struct {
	fakeData fakeDataFunc
}

// New создает объект генератора
func New() Generator {
	return &generator{
		fakeData: faker.FakeData,
	}
}

// GenerateSpecificationData генерирует пример файла спецификации
func (g *generator) GenerateSpecificationData() ([]byte, error) {
	return specExampleSource, nil
}

// GeneratePluginData генерирует данные для .go файла плагина, который возвращает описание структур
func (g *generator) GeneratePluginData(pluginSpec models.PluginDesc) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	if err := specPluginTemplate.Execute(buff, pluginSpec); err != nil {
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
				handlerDesc, err := g.generateProtoHandlerConfig(protoSvc, protoHandler)
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

// GenerateRules генерирует правила из конфигов сервисов
func (g *generator) GenerateRules(svcConfigs []models.ServiceConfigDesc) []*models.Rule {
	var res []*models.Rule
	for _, svcConfig := range svcConfigs {
		for _, handlerCfg := range svcConfig.GRPCHandlers {
			for _, handlerRule := range handlerCfg.Rules {
				r := &models.Rule{
					Conditions: []models.Condition{
						{
							Key:   config.RequestFieldHost,
							Value: svcConfig.Host,
						},
						{
							Key:   config.RequestFieldServiceName,
							Value: handlerCfg.ServiceName,
						},
						{
							Key:   config.RequestFieldMethodName,
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

// GenerateReverseProxyConfigData генерирует конфигурационный файл для прокси сервера
func (g *generator) GenerateReverseProxyConfigData() ([]byte, error) {
	reverseProxyConfig := struct {
		ISOServerHost   string
		HeaderHost      string
		HeaderRequestID string
	}{
		ISOServerHost:   config.ISOServerHost,
		HeaderHost:      config.RequestHeaderHost,
		HeaderRequestID: config.RequestHeaderReqID,
	}
	buff := bytes.NewBuffer(nil)
	if err := reverseProxyConfigTemplate.Execute(buff, reverseProxyConfig); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (g *generator) generateProtoHandlerConfig(protoSvc *public_models.ProtoService, protoHandler public_models.ProtoMethod) (models.HandlerConfigDesc, error) {
	respStruct := protoHandler.ResponseStruct
	if err := g.fakeData(respStruct); err != nil {
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
						Key:   config.RequestFieldReqID,
						Value: "*",
					},
				},
				Response: models.HandlerResponseDesc{
					Data:  string(respData),
					Delay: config.HandlerConfigDefaultTimeout,
				},
			},
		},
	}, nil
}
