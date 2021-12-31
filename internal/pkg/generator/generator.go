package generator

import (
	"bytes"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateConfigData() ([]byte, error) {
	return []byte(configTemplate), nil
}

func (g *Generator) GenerateProtoPluginData(protoPlugin models.ProtoPlugin) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	if err := implTemplate.Execute(buff, protoPlugin); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
