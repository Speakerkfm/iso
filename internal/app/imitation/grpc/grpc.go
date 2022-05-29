package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/metrics"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	public_models "github.com/Speakerkfm/iso/pkg/models"
)

type Handler interface {
	Handle(ctx context.Context, req *Request) (*Response, error)
}

type Response struct {
	msg       proto.Message
	lazyBytes []byte
}

// Service зарегистрированный сервис
type Service struct {
	Methods map[string]Method
}

type Method struct {
	RespStruct proto.Message
}

func New(processor request_processor.Processor, services []*public_models.ProtoService) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ForceServerCodec(codec{}),
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor),
	)
	h := NewHandler(processor)
	for _, svc := range services {
		grpcService := h.registerService(svc)

		srv.RegisterService(grpcService, h)
	}

	return srv
}

func (h *handler) registerService(svc *public_models.ProtoService) *grpc.ServiceDesc {
	grpcMethods := make([]grpc.MethodDesc, 0, len(svc.Methods))
	methods := make(map[string]Method, len(svc.Methods))
	for _, method := range svc.Methods {
		grpcMethods = append(grpcMethods, grpc.MethodDesc{
			Handler:    createUnaryHandler(svc.Name, method.Name, method.RequestStruct),
			MethodName: method.Name,
		})
		methods[method.Name] = Method{
			RespStruct: method.ResponseStruct,
		}
	}
	h.svc[svc.Name] = Service{
		Methods: methods,
	}
	return &grpc.ServiceDesc{
		ServiceName: svc.Name,
		HandlerType: (*Handler)(nil),
		Methods:     grpcMethods,
		Streams:     []grpc.StreamDesc{}, // not implemented. think about it
		Metadata:    svc.ProtoPath,
	}
}

func createUnaryHandler(serviceName, methodName string, msg proto.Message) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := &Request{
			Context:     ctx,
			ServiceName: serviceName,
			MethodName:  methodName,
			Msg:         msg.ProtoReflect().New().Interface(),
			HandledAt:   time.Now(),
		}
		if err := dec(in); err != nil {
			logger.Errorf(ctx, "fail to decode incoming req: %s", err.Error())
			return nil, err
		}

		if interceptor == nil {
			return srv.(Handler).Handle(ctx, in)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: fmt.Sprintf("/%s/%s", serviceName, methodName),
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(Handler).Handle(ctx, req.(*Request))
		}
		return interceptor(ctx, in, info, handler)
	}
}
