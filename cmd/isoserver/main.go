package main

import (
	"log"
	"net"
	"plugin"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	"github.com/Speakerkfm/iso/internal/pkg/rule/manager"
	models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	pluginPath = "struct.so"
	serverHost = "localhost:8001"
)

func main() {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Fatalf("fail to open plugin: %s", pluginPath)
	}

	svcs, err := plug.Lookup(models.ServiceProviderName)
	if err != nil {
		log.Fatalf("fail too look up ServiceProvider in plugin: %s", err.Error())
	}

	s, ok := svcs.(models.ServiceProvider)
	if !ok {
		log.Fatal("fail to get proto description from module")
	}

	lis, err := net.Listen("tcp", serverHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ruleManager := manager.New()
	processor := request_processor.New(ruleManager)

	impl := imitation.New(processor, s.GetList())

	log.Println("iso server created")

	if err := impl.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
