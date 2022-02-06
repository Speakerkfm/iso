package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"plugin"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/pkg/util"
	models "github.com/Speakerkfm/iso/pkg/models"
)

type serviceProvider interface {
	GetList() []*models.ProtoService
}

func main() {
	plug, err := plugin.Open("struct.so")
	if err != nil {
		panic(err)
	}

	svcs, err := plug.Lookup("ServiceProvider")
	if err != nil {
		panic(err)
	}

	s, ok := svcs.(serviceProvider)
	if !ok {
		fmt.Printf("convert failed")
	}
	lis, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.CustomCodec(codec{}))
	grpcServer := grpc.NewServer(opts...)
	i := NewImpl(s.GetList())

	for _, svc := range s.GetList() {
		grpcService := createGRPCService(svc)

		grpcServer.RegisterService(grpcService, i)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type codec struct {
}

func (codec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(*Response)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	if vv.lazyBytes != nil {
		return vv.lazyBytes, nil
	}
	return proto.Marshal(vv.msg)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(*Request)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv.Msg)
}

func (codec) String() string {
	return "kek"
}

type Request struct {
	ServiceName string
	MethodName  string
	Msg         proto.Message
}

type Response struct {
	msg       proto.Message
	lazyBytes []byte
}

type impl struct {
	storedResponses map[string]*Response
	rules           map[string]map[string][]*Rule
}

func NewImpl(services []*models.ProtoService) *impl {
	storedResponses := make(map[string]*Response)
	rules := make(map[string]map[string][]*Rule)
	protoStructs := make(map[string]map[string]models.ProtoMethod)
	for _, svc := range services {
		if protoStructs[svc.Name] == nil {
			protoStructs[svc.Name] = make(map[string]models.ProtoMethod)
		}
		for _, m := range svc.Methods {
			protoStructs[svc.Name][m.Name] = m
		}
	}

	for _, r := range getRawRules() {
		if rules[r.ServiceName] == nil {
			rules[r.ServiceName] = make(map[string][]*Rule)
		}

		respID := util.NewUUID()

		rules[r.ServiceName][r.MethodName] = append(rules[r.ServiceName][r.MethodName], &Rule{
			respID: respID,
		})

		msg := protoStructs[r.ServiceName][r.MethodName].ResponseStruct.ProtoReflect().New().Interface()

		if err := json.Unmarshal(r.ResponseData, msg); err != nil {
			log.Fatalf("failed to unmarshal: %v", err)
		}

		lazyBytes, err := proto.Marshal(msg)
		if err != nil {
			log.Fatalf("failed to marshal lazy bytes: %v", err)
		}

		storedResponses[respID] = &Response{
			msg:       msg,
			lazyBytes: lazyBytes,
		}
	}

	return &impl{
		storedResponses: storedResponses,
		rules:           rules,
	}
}

type Rule struct {
	conditions []condition
	respID     string
}

func (r *Rule) IsSelected() bool {
	return true
}

type condition struct {
	Key   string
	Value interface{}
}

type RawRule struct {
	ServiceName  string
	MethodName   string
	Conditions   []condition
	ResponseData json.RawMessage
}

func getRawRules() []*RawRule {
	return []*RawRule{
		{
			ServiceName: "UserService",
			MethodName:  "GetUser",
			Conditions: []condition{
				{
					Key:   "id",
					Value: 10,
				},
			},
			ResponseData: []byte(`{"user":{"id":10,"name":"kek"}}`),
		},
		{
			ServiceName: "PhoneService",
			MethodName:  "CheckPhone",
			Conditions: []condition{
				{
					Key:   "id",
					Value: 10,
				},
			},
			ResponseData: []byte(`{"exists":true}`),
		},
	}
}

func (i *impl) Process(ctx context.Context, req *Request) (*Response, error) {
	fmt.Printf("message: %+v\n", req.Msg.ProtoReflect().Descriptor().Fields())

	rules := i.rules[req.ServiceName][req.MethodName]

	resp := i.storedResponses[rules[0].respID]

	return resp, nil
}

func createGRPCService(svc *models.ProtoService) *grpc.ServiceDesc {
	grpcMethods := make([]grpc.MethodDesc, 0, len(svc.Methods))
	for _, method := range svc.Methods {
		grpcMethods = append(grpcMethods, grpc.MethodDesc{
			Handler:    requestHandler(svc.Name, method.Name, method.RequestStruct),
			MethodName: method.Name,
		})
	}
	return &grpc.ServiceDesc{
		ServiceName: svc.Name,
		HandlerType: (*RequestProcessor)(nil),
		Methods:     grpcMethods,
		Streams:     []grpc.StreamDesc{},
		Metadata:    svc.ProtoPath,
	}
}

func requestHandler(serviceName, methodName string, msg proto.Message) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := &Request{
			ServiceName: serviceName,
			MethodName:  methodName,
			Msg:         msg.ProtoReflect().New().Interface(),
		}
		if err := dec(in); err != nil {
			fmt.Printf("err: %+v\n", err)
			return nil, err
		}
		fmt.Printf("in: %+v\n", in)
		if interceptor == nil {
			return srv.(RequestProcessor).Process(ctx, in)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: fmt.Sprintf("/%s/%s", serviceName, methodName),
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(RequestProcessor).Process(ctx, req.(*Request))
		}
		return interceptor(ctx, in, info, handler)
	}
}

type RequestProcessor interface {
	Process(ctx context.Context, req *Request) (*Response, error)
}
