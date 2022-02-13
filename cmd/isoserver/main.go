package main

import (
	"context"
	"log"
	"net"
	"plugin"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	"github.com/Speakerkfm/iso/internal/pkg/rule/manager"
	rule_parser "github.com/Speakerkfm/iso/internal/pkg/rule/parser"
	models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	pluginPath = "struct.so"      // in args ...
	serverHost = "localhost:8001" // in args ...
)

func main() {
	appCtx := context.Background()

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

	ruleParser := rule_parser.New()
	rules, err := ruleParser.Parse(appCtx, "")
	if err != nil {
		log.Fatalf("fail to parse rules: %s", err.Error())
	}

	ruleManager := manager.New()
	ruleManager.UpdateRuleTree(rules)

	processor := request_processor.New(ruleManager)

	impl := imitation.New(processor, s.GetList())

	log.Println("iso server created")

	if err := impl.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
