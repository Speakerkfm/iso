package proto_parser

import (
	"fmt"
	"strings"

	"github.com/emicklei/proto"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(rawProtoFile []byte) ([]*models.ProtoServiceDesc, error) {
	fileReader := strings.NewReader(string(rawProtoFile))
	protoParser := proto.NewParser(fileReader)

	definition, err := protoParser.Parse()
	if err != nil {
		return nil, err
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

func handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
}

func handleRPC(r *proto.RPC) {
	fmt.Println(r.RequestType)
	fmt.Println(r.ReturnsType)
}
