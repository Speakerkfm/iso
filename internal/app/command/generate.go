package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/util"
	public_models "github.com/Speakerkfm/iso/pkg/models"
)

const (
	pluginGoFile     = "spec.go"
	pluginModuleName = "iso_proto"
)

func (c *Command) Generate(ctx context.Context, specPath string) error {
	cmdExecDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to get cmd exec dir: %w", err)
	}
	currDir := filepath.Join(cmdExecDir, filepath.Dir(specPath))
	logger.Infof(ctx, "Current directory: %s", currDir)

	logger.Info(ctx, "Loading specification...")
	spec, err := c.loadSpec(specPath)
	if err != nil {
		return fmt.Errorf("fail to load specification from path %s: %w", specPath, err)
	}

	logger.Info(ctx, "Processing specification...")
	protoFiles, err := c.processSpec(spec)
	if err != nil {
		return fmt.Errorf("fail to process config: %w", err)
	}
	logger.Infof(ctx, "Found %d proto files", len(protoFiles))

	wd := fmt.Sprintf("%s%s", os.TempDir(), util.NewUUID())
	if err := os.MkdirAll(wd, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to make temp dir %s: %w", wd, err)
	}
	logger.Infof(ctx, "Temp directory: %s", wd)

	logger.Info(ctx, "Processing spec files...")
	pluginSpec, err := c.processSpecFiles(ctx, wd, protoFiles)
	if err != nil {
		return fmt.Errorf("fail to process spec files: %w", err)
	}

	logger.Info(ctx, "Generating plugin data...")
	pluginData, err := c.gen.GeneratePluginData(pluginSpec)
	if err != nil {
		return fmt.Errorf("fail to generate proto plugin: %w", err)
	}

	if err := ioutil.WriteFile(filepath.Join(wd, pluginGoFile), pluginData, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to save plugin data to file: %w", err)
	}
	logger.Info(ctx, "Plugin data generated")

	logger.Info(ctx, "Building spec plugin...")
	if err := c.golang.BuildPlugin(wd, currDir, pluginModuleName, pluginGoFile); err != nil {
		return fmt.Errorf("fail to build proto plugin: %w", err)
	}
	logger.Info(ctx, "Proto plugin generated")

	logger.Info(ctx, "Generating rule examples....")
	if err := c.generateRuleExamples(ctx, spec, currDir); err != nil {
		return fmt.Errorf("fail to generate rules example: %w", err)
	}

	logger.Info(ctx, "Done")
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

func (c *Command) processSpecFiles(ctx context.Context, wd string, protoFiles []*models.ProtoFileData) (models.PluginDesc, error) {
	var err error
	protoPlugin := models.PluginDesc{
		ModuleName: pluginModuleName,
	}

	processedServices := make(map[string]struct{})

	for _, protoFile := range protoFiles {
		logger.Infof(ctx, "Fetching %s file...", protoFile.OriginalPath)
		protoFile.RawData, err = c.fileFetcher.FetchFile(ctx, protoFile.OriginalPath)
		if err != nil {
			return models.PluginDesc{}, err
		}

		protoDir := filepath.Join(wd, protoFile.PkgName)
		if err := os.MkdirAll(protoDir, fs.ModePerm); err != nil {
			return models.PluginDesc{}, err
		}

		if err := ioutil.WriteFile(filepath.Join(wd, protoFile.Path), protoFile.RawData, fs.ModePerm); err != nil {
			return models.PluginDesc{}, err
		}

		logger.Infof(ctx, "Processing %s file with protoc...", protoFile.OriginalPath)
		if err := c.protoc.Process(wd, protoFile); err != nil {
			return models.PluginDesc{}, err
		}

		serviceDescriptions, err := c.protoParser.Parse(protoFile.RawData)
		if err != nil {
			return models.PluginDesc{}, err
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

		logger.Infof(ctx, "File %s was successfully processed", protoFile.OriginalPath)
	}

	return protoPlugin, nil
}

func (c *Command) generateRuleExamples(ctx context.Context, spec models.ServiceSpecification, path string) error {
	svcProvider, err := loadPluginData(ctx, filepath.Join(path, models.PluginName))
	if err != nil {
		return fmt.Errorf("fail to load plugin data: %w", err)
	}

	svcConfigs, err := c.gen.GenerateServiceConfigs(spec, svcProvider)
	if err != nil {
		return fmt.Errorf("fail to generate svc cfgs: %w", err)
	}

	currDir := filepath.Join(path, "rules")
	if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
	}

	for _, svcCfg := range svcConfigs {
		currDir := filepath.Join(currDir, svcCfg.Name)
		if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
			return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
		}
		if err := encodeDataToYamlFile(svcCfg, currDir, "service.yaml"); err != nil {
			return fmt.Errorf("fail to encode data to yaml: %w", err)
		}
		if len(svcCfg.GRPCHandlers) > 0 {
			currDir := filepath.Join(currDir, "proto")
			if err := os.Mkdir(currDir, fs.ModePerm); err != nil {
				return fmt.Errorf("fail to create dir: %s, err: %w", currDir, err)
			}
			for _, handlerCfg := range svcCfg.GRPCHandlers {
				fileName := fmt.Sprintf("%s_%s.yaml", handlerCfg.ServiceName, handlerCfg.MethodName)
				if err := encodeDataToYamlFile(handlerCfg, currDir, fileName); err != nil {
					return fmt.Errorf("fail to encode data to yaml: %w", err)
				}
			}
		}
	}

	return nil
}

func encodeDataToYamlFile(data interface{}, currDir, fileName string) error {
	f, err := os.Create(filepath.Join(currDir, fileName))
	if err != nil {
		return fmt.Errorf("fail to create %s: %w", fileName, err)
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("fail to encode data: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("fail to close yaml file: %w", err)
	}
	return nil
}

func loadPluginData(ctx context.Context, pluginPath string) (public_models.ServiceProvider, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("fail to open plugin: %s, err: %w", pluginPath, err)
	}

	svcs, err := plug.Lookup(public_models.ServiceProviderName)
	if err != nil {
		return nil, fmt.Errorf("fail too look up ServiceProvider in plugin: %s", err.Error())
	}

	s, ok := svcs.(public_models.ServiceProvider)
	if !ok {
		return nil, fmt.Errorf("fail to get proto description from module")
	}

	return s, nil
}
