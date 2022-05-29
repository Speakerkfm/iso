package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"plugin"
	"syscall"

	"go.uber.org/atomic"
	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"
	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/app/admin"
	"github.com/Speakerkfm/iso/internal/app/imitation"
	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/events"
	events_batcher "github.com/Speakerkfm/iso/internal/pkg/events/batcher"
	event_store "github.com/Speakerkfm/iso/internal/pkg/events/repository/store"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/metrics"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/reporter"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
	"github.com/Speakerkfm/iso/internal/pkg/router"
	"github.com/Speakerkfm/iso/internal/pkg/rule/manager"
	rule_parser "github.com/Speakerkfm/iso/internal/pkg/rule/parser"
	"github.com/Speakerkfm/iso/internal/pkg/rule/repository/store"
	"github.com/Speakerkfm/iso/internal/pkg/rule/service"
	"github.com/Speakerkfm/iso/internal/pkg/rule/syncer"
	public_models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	configPath = "" // now config with constants
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())

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

	st := store.New()
	if err := st.SaveServiceConfigs(appCtx, serviceConfigs); err != nil {
		logger.Fatalf(appCtx, "fail to save service configs: %s", err.Error())
	}

	ruleSvc := service.New(st, gen)

	defaultRules, err := ruleSvc.GetRules(appCtx)
	if err != nil {
		logger.Fatalf(appCtx, "fail to get default rules: %s", err.Error())
	}

	ruleManager := manager.New()
	ruleManager.UpdateRuleTree(defaultRules)

	eventRepo := event_store.New()

	batcher := events_batcher.New(appCtx,
		atomic.NewBool(config.BatcherEnabled),
		eventRepo,
		config.BatcherBatchCount,
		atomic.NewDuration(config.BatcherFlushInterval),
		atomic.NewInt64(config.BatcherFlushItemsAmount),
		config.BatcherEventBuffSize,
	)
	eventSvc := events.New(batcher, eventRepo)
	processor := request_processor.New(ruleManager, eventSvc)

	ruleSyncer := syncer.New(appCtx, ruleSvc, ruleManager, config.RulesSyncInterval)
	ruleSyncer.Start()
	defer ruleSyncer.Stop()

	reportSvc := reporter.New(eventSvc)

	adminServer := admin.New(ruleSvc, reportSvc)

	impl := imitation.New(processor, s.GetList())

	done := make(chan struct{})
	go func() {
		lis, err := net.Listen("tcp", config.ISOServerGRPCHost)
		if err != nil {
			logger.Fatalf(appCtx, "failed to listen: %v", err)
		}
		defer lis.Close()
		if err := impl.RegisterGRPC(lis); err != nil {
			logger.Fatalf(appCtx, "failed to serve: %s", err.Error())
		}
		close(done)
	}()
	go func() {
		mx := router.NewRouter()
		if err := metrics.RegisterMetricsHandler(appCtx, mx); err != nil {
			logger.Fatalf(appCtx, "fail to register metrics handler: %s", err.Error())
		}
		if err := adminServer.RegisterGateway(appCtx, mx); err != nil {
			logger.Fatalf(appCtx, "fail to register admin gateway: %s", err.Error())
		}
		if err := http.ListenAndServe(config.ISOServerAdminHost, mx); err != nil {
			logger.Fatalf(appCtx, "failed to serve: %s", err.Error())
		}
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		<-c
		cancel()
		close(done)
	}()

	logger.Info(appCtx, "iso server created")
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
