package imitation

import (
	"net"

	"google.golang.org/grpc"

	grpc_server "github.com/Speakerkfm/iso/internal/app/imitation/grpc"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	public_models "github.com/Speakerkfm/iso/pkg/models"
)

type Imitation struct {
	grpcServer *grpc.Server
}

func New(processor request_processor.Processor, protoServices []*public_models.ProtoService) *Imitation {
	grpcServer := grpc_server.New(processor, protoServices)

	return &Imitation{
		grpcServer: grpcServer,
	}
}

func (i *Imitation) RegisterGRPC(lis net.Listener) error {
	return i.grpcServer.Serve(lis)
}
