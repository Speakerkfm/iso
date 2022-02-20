package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"github.com/bxcodec/faker/v3"
	"gopkg.in/yaml.v2"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/util"
	shared_models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	protoPluginFile = "struct.go"
	protoPluginName = "iso_proto"
)

func (c *Command) Generate(ctx context.Context, specPath string) error {
	logger.Info(ctx, "Loading specification...")
	spec, err := c.loadSpec(specPath)
	if err != nil {
		return fmt.Errorf("fail to load specification from path %s: %w", specPath, err)
	}

	protoFiles, err := c.processSpec(spec)
	if err != nil {
		return fmt.Errorf("fail to process config: %w", err)
	}

	wd := fmt.Sprintf("%s%s", os.TempDir(), util.NewUUID())

	logger.Infof(ctx, "Working directory: %s\n", wd)

	if err := os.MkdirAll(wd, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to make temp dir %s: %w", wd, err)
	}

	logger.Info(ctx, "Processing proto files...")

	protoPlugin, err := c.processProtoFiles(ctx, wd, protoFiles)
	if err != nil {
		return fmt.Errorf("fail to process proto files: %w", err)
	}

	protoPluginData, err := c.gen.GenerateProtoPluginData(protoPlugin)
	if err != nil {
		return fmt.Errorf("fail to generate proto plugin: %w", err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", wd, protoPluginFile), protoPluginData, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to save proto plugin data to file: %w", err)
	}

	if err := c.golang.CreateModule(wd, protoPluginName); err != nil {
		return fmt.Errorf("fail to create go mod for plugin: %w", err)
	}

	logger.Info(ctx, "Building proto plugin...")
	currDir, _ := os.Getwd()
	outDir := filepath.Join(currDir, filepath.Dir(specPath))
	logger.Infof(ctx, "out dir: %s", outDir)
	if err := c.golang.BuildPlugin(wd, outDir, protoPluginFile); err != nil {
		return fmt.Errorf("fail to build proto plugin: %w", err)
	}

	logger.Info(ctx, "Proto plugin generated")

	if err := c.generateRulesExample(ctx, spec, outDir); err != nil {
		return fmt.Errorf("fail to generate rules example: %w", err)
	}

	return nil
}

func (c *Command) loadSpec(path string) (models.ServiceSpecification, error) {
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

func (c *Command) processSpec(spec models.ServiceSpecification) ([]*models.ProtoFileData, error) {
	var protoFiles []*models.ProtoFileData

	for _, dep := range spec.ExternalDependencies {
		for _, protoPath := range dep.ProtoPaths {
			fileName := filepath.Base(protoPath)
			pkgName := strings.Split(fileName, ".")[0]
			protoFiles = append(protoFiles, &models.ProtoFileData{
				Name:         fileName,
				PkgName:      pkgName,
				OriginalPath: protoPath,
				Path:         fmt.Sprintf("%s/%s", pkgName, fileName),
			})
		}
	}

	return protoFiles, nil
}

func (c *Command) processProtoFiles(ctx context.Context, wd string, protoFiles []*models.ProtoFileData) (models.ProtoPluginDesc, error) {
	var err error
	protoPlugin := models.ProtoPluginDesc{
		ModuleName: protoPluginName,
	}

	processedServices := make(map[string]struct{})

	for _, protoFile := range protoFiles {
		protoFile.RawData, err = c.fileFetcher.FetchFile(ctx, protoFile.OriginalPath)
		if err != nil {
			return models.ProtoPluginDesc{}, err
		}

		protoDir := fmt.Sprintf("%s/%s", wd, protoFile.PkgName)
		if err := os.MkdirAll(protoDir, fs.ModePerm); err != nil {
			return models.ProtoPluginDesc{}, err
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", wd, protoFile.Path), protoFile.RawData, fs.ModePerm); err != nil {
			return models.ProtoPluginDesc{}, err
		}

		if err := c.protoc.Process(wd, protoFile); err != nil {
			return models.ProtoPluginDesc{}, err
		}

		serviceDescriptions, err := c.protoParser.Parse(protoFile.RawData)
		if err != nil {
			return models.ProtoPluginDesc{}, err
		}

		for _, svcDesc := range serviceDescriptions {
			if _, isProcessed := processedServices[svcDesc.Name]; isProcessed {
				continue
			}

			svcDesc.ProtoPath = protoFile.Path
			svcDesc.PkgName = protoFile.PkgName

			protoPlugin.Imports = append(protoPlugin.Imports, fmt.Sprintf("%s \"%s/%s\"", svcDesc.PkgName, protoPlugin.ModuleName, svcDesc.PkgName))
			protoPlugin.ProtoServices = append(protoPlugin.ProtoServices, svcDesc)

			processedServices[svcDesc.Name] = struct{}{}
		}
	}

	return protoPlugin, nil
}

func (c *Command) generateRulesExample(ctx context.Context, spec models.ServiceSpecification, path string) error {
	svcProvider, err := loadPluginData(ctx, filepath.Join(path, models.PluginName))
	if err != nil {
		return fmt.Errorf("fail to load plugin data: %w", err)
	}

	dataPathProtoService := make(map[string]*shared_models.ProtoService)

	for _, svc := range svcProvider.GetList() {
		dataPathProtoService[svc.ProtoPath] = svc
	}

	var svcExamples []models.ServiceConfigDesc
	{
	}

	for _, dep := range spec.ExternalDependencies {
		svcDesc := models.ServiceConfigDesc{
			Host: dep.Host,
		}
		for _, protoDepPath := range dep.ProtoPaths {
			protoSvc, ok := dataPathProtoService[protoDepPath]
			if !ok {
				// todo
			}

			for _, protoHandler := range protoSvc.Methods {
				respStruct := protoHandler.ResponseStruct
				if err := faker.FakeData(respStruct); err != nil {

				}
				respData, err := json.MarshalIndent(respStruct, "", "	")
				if err != nil {

				}
				handlerDesc := models.HandlerConfigDesc{
					ServiceName: protoSvc.Name,
					MethodName:  protoHandler.Name,
					Rules: []models.RuleDesc{
						// generate
						{
							Conditions: []models.HandlerConditionDesc{
								{
									Key:   "header.x-request-id",
									Value: "*",
								},
							},
							Response: models.HandlerResponseDesc{
								Data:  string(respData),
								Delay: 5 * time.Millisecond,
							},
						},
					},
				}
				svcDesc.GRPCHandlers = append(svcDesc.GRPCHandlers, handlerDesc)
			}
		}
		svcExamples = append(svcExamples, svcDesc)
	}

	fmt.Printf("generated rules: %+v\n", svcExamples)
	return nil
}

func loadPluginData(ctx context.Context, pluginPath string) (shared_models.ServiceProvider, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("fail to open plugin: %s, err: %w", pluginPath, err)
	}

	svcs, err := plug.Lookup(shared_models.ServiceProviderName)
	if err != nil {
		return nil, fmt.Errorf("fail too look up ServiceProvider in plugin: %s", err.Error())
	}

	s, ok := svcs.(shared_models.ServiceProvider)
	if !ok {
		return nil, fmt.Errorf("fail to get proto description from module")
	}

	return s, nil
}
