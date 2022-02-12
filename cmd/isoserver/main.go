package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"plugin"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	models "github.com/Speakerkfm/iso/pkg/models"
)

type ServiceProvider interface {
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

	s, ok := svcs.(ServiceProvider)
	if !ok {
		fmt.Printf("convert failed")
	}

	lis, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	processor := request_processor.New()

	impl := imitation.New(processor, s.GetList())
	if err := impl.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
	Value string
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
					Value: "10",
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
					Value: "10",
				},
			},
			ResponseData: []byte(`{"exists":true}`),
		},
	}
}
