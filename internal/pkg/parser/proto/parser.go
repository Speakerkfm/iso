package proto

import (
	"bytes"
	"fmt"

	"github.com/emicklei/proto"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Parser парсит .proto файлы
type Parser interface {
	Parse(rawProtoFile []byte) ([]*models.ProtoServiceDesc, error)
}

type parser struct {
}

// New создает новый парсер
func New() Parser {
	return &parser{}
}

// Parse парсит пришедший .proto файл во внутреннюю структуру, которая содержит описание proto структур
func (p *parser) Parse(rawProtoFile []byte) ([]*models.ProtoServiceDesc, error) {
	protoParser := proto.NewParser(bytes.NewBuffer(rawProtoFile))

	definition, err := protoParser.Parse()
	if err != nil {
		return nil, fmt.Errorf("fail to parse proto file: %w", err)
	}

	resBuff := &resultBuffer{}

	proto.Walk(definition, proto.WithService(resBuff.handleService))

	return resBuff.serviceDescriptions, nil
}

type resultBuffer struct {
	serviceDescriptions []*models.ProtoServiceDesc
}

func (rb *resultBuffer) handleService(s *proto.Service) {
	svcDesc := &models.ProtoServiceDesc{
		Name: s.Name,
	}
	for _, e := range s.Elements {
		if rpc, ok := e.(*proto.RPC); ok {
			svcDesc.Methods = append(svcDesc.Methods, &models.ProtoMethodDesc{
				Name:         rpc.Name,
				RequestType:  rpc.RequestType,
				ResponseType: rpc.ReturnsType,
			})
		}
	}

	rb.serviceDescriptions = append(rb.serviceDescriptions, svcDesc)
}
