package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"path"
	"plugin"
	"syscall"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"
	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	"github.com/Speakerkfm/iso/internal/pkg/rule/manager"
	rule_parser "github.com/Speakerkfm/iso/internal/pkg/rule/parser"
	public_models "github.com/Speakerkfm/iso/pkg/models"
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
	specPath := path.Join(*isoDir, config.SpecificationFileName)

	if err := config.Parse(configPath); err != nil {
		logger.Fatalf(appCtx, "fail to parse config file: %s, err: %w", configPath, err)
	}

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		logger.Fatalf(appCtx, "fail to open plugin: %s, err: %s", pluginPath, err.Error())
	}

	svcs, err := plug.Lookup(public_models.ServiceProviderName)
	if err != nil {
		logger.Fatalf(appCtx, "fail too look up ServiceProvider in plugin: %s", err.Error())
	}

	s, ok := svcs.(public_models.ServiceProvider)
	if !ok {
		logger.Fatalf(appCtx, "fail to get proto description from module")
	}

	gen := generator.New()

	var serviceConfigs []models.ServiceConfigDesc
	if _, err := os.Stat(ruleDirectoryPath); os.IsNotExist(err) {
		spec, err := loadSpec(specPath)
		if err != nil {
			logger.Fatalf(appCtx, "fail to load spec: %s", err.Error())
		}

		serviceConfigs, err = gen.GenerateServiceConfigs(spec, s)
		if err != nil {
			logger.Fatalf(appCtx, "fail to generate service configs: %s", err.Error())
		}
	} else {
		ruleParser := rule_parser.New()
		serviceConfigs, err = ruleParser.ParseDirectory(appCtx, ruleDirectoryPath)
		if err != nil {
			logger.Fatalf(appCtx, "fail to parse rules: %s", err.Error())
		}
	}

	rules := gen.GenerateRules(serviceConfigs)

	ruleManager := manager.New()
	ruleManager.UpdateRuleTree(rules)

	processor := request_processor.New(ruleManager)

	logger.Info(appCtx, "iso server created")

	impl := imitation.New(processor, s.GetList())

	lis, err := net.Listen("tcp", config.ISOServerHost)
	if err != nil {
		logger.Fatalf(appCtx, "failed to listen: %v", err)
	}
	defer lis.Close()

	done := make(chan struct{})
	go func() {
		if err := impl.Serve(lis); err != nil {
			logger.Fatalf(appCtx, "failed to serve: %v", err)
		}
		close(done)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		<-c
		close(done)
	}()

	<-done
	logger.Infof(appCtx, "iso server stopped")
}

func loadSpec(path string) (models.ServiceSpecification, error) {
	spec := models.ServiceSpecification{}

	fin, err := os.Open(path)
	if err != nil {
		return models.ServiceSpecification{}, err
	}
	defer fin.Close()

	if err := yaml.NewDecoder(fin).Decode(&spec); err != nil {
		return models.ServiceSpecification{}, err
	}

	return spec, nil
}
