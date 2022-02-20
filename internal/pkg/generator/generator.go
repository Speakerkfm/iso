package generator

import (
	"bytes"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Generator генерирует файлы по шаблонам
type Generator interface {
	GenerateSpecificationData() ([]byte, error)
	GenerateProtoPluginData(protoPlugin models.ProtoPluginDesc) ([]byte, error)
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

// GenerateProtoPluginData генерирует .go файл прото плагина, который возвращает описание прото структур
func (g *generator) GenerateProtoPluginData(protoPlugin models.ProtoPluginDesc) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	if err := implTemplate.Execute(buff, protoPlugin); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
