package generator

import (
	"bytes"
	"fmt"

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

func (g *Generator) GenerateProtoPlugin(moduleName string, svcDesc map[string]*models.ProtoServiceDesc) ([]byte, error) {
	protoPlugin := models.ProtoPlugin{
		ModuleName: moduleName,
	}
	for _, svc := range svcDesc {
		protoPlugin.Imports = append(protoPlugin.Imports, fmt.Sprintf("%s \"%s/pb/%s\"", svc.PkgName, protoPlugin.ModuleName, svc.PkgName))
		protoPlugin.ProtoServices = append(protoPlugin.ProtoServices, svc)
	}

	buff := bytes.NewBuffer(nil)
	if err := implTemplate.Execute(buff, protoPlugin); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
