package main

import (
	"context"
	"flag"
	"net"
	"path"
	"plugin"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	"github.com/Speakerkfm/iso/internal/pkg/rule/manager"
	rule_parser "github.com/Speakerkfm/iso/internal/pkg/rule/parser"
	shared_models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	configPath = "" // now config with constants
)

func main() {
	appCtx := context.Background()

	isoDir := flag.String("dir", config.DefaultProjectDir, "directory with rules and plugin")
	flag.Parse()

	pluginPath := path.Join(*isoDir, config.PluginFileName)
	ruleDirectoryPath := path.Join(*isoDir, config.RulesDir)

	if err := config.Parse(configPath); err != nil {
		logger.Fatalf(appCtx, "fail to parse config file: %s, err: %w", configPath, err)
	}

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		logger.Fatalf(appCtx, "fail to open plugin: %s, err: %w", pluginPath, err)
	}

	svcs, err := plug.Lookup(shared_models.ServiceProviderName)
	if err != nil {
		logger.Fatalf(appCtx, "fail too look up ServiceProvider in plugin: %s", err.Error())
	}

	s, ok := svcs.(shared_models.ServiceProvider)
	if !ok {
		logger.Fatalf(appCtx, "fail to get proto description from module")
	}

	lis, err := net.Listen("tcp", config.ISOServerHost)
	if err != nil {
		logger.Fatalf(appCtx, "failed to listen: %v", err)
	}

	ruleParser := rule_parser.New()
	serviceConfigs, err := ruleParser.ParseDirectory(appCtx, ruleDirectoryPath)
	if err != nil {
		logger.Fatalf(appCtx, "fail to parse rules: %s", err.Error())
	}

	gen := generator.New()
	rules := gen.GenerateRules(serviceConfigs)

	ruleManager := manager.New()
	ruleManager.UpdateRuleTree(rules)

	processor := request_processor.New(ruleManager)

	impl := imitation.New(processor, s.GetList())

	logger.Info(appCtx, "iso server created")

	if err := impl.Serve(lis); err != nil {
		logger.Fatalf(appCtx, "failed to serve: %v", err)
	}
}
