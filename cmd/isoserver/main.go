package main

import (
	"context"
	"net"
	"plugin"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	"github.com/Speakerkfm/iso/internal/pkg/rule/manager"
	rule_parser "github.com/Speakerkfm/iso/internal/pkg/rule/parser"
	models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	pluginPath        = "struct.so"       // in args ...
	ruleDirectoryPath = "./example/rules" // in args ...
	serverHost        = "localhost:8001"  // in args ...
)

func main() {
	appCtx := context.Background()

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		logger.Fatalf(appCtx, "fail to open plugin: %s", pluginPath)
	}

	svcs, err := plug.Lookup(models.ServiceProviderName)
	if err != nil {
		logger.Fatalf(appCtx, "fail too look up ServiceProvider in plugin: %s", err.Error())
	}

	s, ok := svcs.(models.ServiceProvider)
	if !ok {
		logger.Fatalf(appCtx, "fail to get proto description from module")
	}

	lis, err := net.Listen("tcp", serverHost)
	if err != nil {
		logger.Fatalf(appCtx, "failed to listen: %v", err)
	}

	ruleParser := rule_parser.New()
	serviceConfigs, err := ruleParser.ParseDirectory(appCtx, ruleDirectoryPath)
	if err != nil {
		logger.Fatalf(appCtx, "fail to parse rules: %s", err.Error())
	}

	rules := ruleParser.GenerateRules(serviceConfigs)

	ruleManager := manager.New()
	ruleManager.UpdateRuleTree(rules)

	processor := request_processor.New(ruleManager)

	impl := imitation.New(processor, s.GetList())

	logger.Info(appCtx, "iso server created")

	if err := impl.Serve(lis); err != nil {
		logger.Fatalf(appCtx, "failed to serve: %v", err)
	}
}
