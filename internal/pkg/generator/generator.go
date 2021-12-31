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

func (g *Generator) GenerateConfig() ([]byte, error) {
	return []byte(configTemplate), nil
}

func (g *Generator) GenerateProtoPlugin(svcDesc map[string]*models.ProtoServiceDesc) ([]byte, error) {
	protoPlugin := models.ProtoPlugin{
		ModuleName: "iso/proto",
	}
	for _, svc := range svcDesc {
		protoPlugin.ProtoServices = append(protoPlugin.ProtoServices, svc)
	}

	buff := bytes.NewBuffer(nil)
	if err := implTemplate.Execute(buff, protoPlugin); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
