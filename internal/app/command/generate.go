package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// Generate - генерирует все необходимые данные (список правил, конфиг прокси и плагин) для проекта в указанной директории
func (c *Command) Generate(ctx context.Context, dir string, dockerEnabled bool) error {
	cmdExecDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to get cmd exec dir: %w", err)
	}
	projectFullDir := filepath.Join(cmdExecDir, dir)
	logger.Infof(ctx, "Current project directory: %s", projectFullDir)

	logger.Info(ctx, "Loading specification...")
	specPath := filepath.Join(projectFullDir, config.SpecificationFileName)
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

	pluginDir := filepath.Join(projectFullDir, config.PluginDir)
	if err := os.MkdirAll(pluginDir, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to make plugin dir %s: %w", pluginDir, err)
	}

	logger.Info(ctx, "Processing spec files...")
	pluginSpec, err := c.processSpecFiles(ctx, pluginDir, protoFiles)
	if err != nil {
		return fmt.Errorf("fail to process spec files: %w", err)
	}

	logger.Info(ctx, "Generating data for plugin...")
	pluginData, err := c.gen.GeneratePluginData(pluginSpec)
	if err != nil {
		return fmt.Errorf("fail to generate data for plugin: %w", err)
	}

	if err := ioutil.WriteFile(filepath.Join(pluginDir, config.PluginGoFileName), pluginData, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to save data for plugin to file: %w", err)
	}
	logger.Info(ctx, "Data for plugin was generated")

	logger.Info(ctx, "Building spec plugin...")
	if err := c.buildPlugin(ctx, dockerEnabled, pluginDir, projectFullDir, protoFiles); err != nil {
		return fmt.Errorf("fail to build plugin: %w", err)
	}
	logger.Info(ctx, "Plugin was generated")

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

func (c *Command) buildPlugin(ctx context.Context, dockerEnabled bool, pluginDir, projectFullDir string, protoFiles []*models.ProtoFileData) error {
	if dockerEnabled {
		if err := c.docker.BuildPlugin(pluginDir, projectFullDir, config.PluginModuleName, config.PluginGoFileName); err != nil {
			return err
		}
		return nil
	}

	for _, protoFile := range protoFiles {
		logger.Infof(ctx, "Processing %s file with protoc...", protoFile.OriginalPath)
		if err := c.protoc.Process(pluginDir, protoFile); err != nil {
			return err
		}
	}
	if err := c.golang.BuildPlugin(pluginDir, filepath.Join(projectFullDir, config.PluginFileName), config.PluginModuleName, config.PluginGoFileName); err != nil {
		return err
	}

	logger.Info(ctx, "Cleaning...")
	if err := os.RemoveAll(pluginDir); err != nil {
		return fmt.Errorf("fail to clean plugin dir: %w", err)
	}
	return nil
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
				Path:         fmt.Sprintf("%s/%s", pkgName, fileName), // TODO подумать: один прото файл у нескольких сервисов
			})
		}
	}

	return protoFiles, nil
}

func (c *Command) processSpecFiles(ctx context.Context, pluginDir string, protoFiles []*models.ProtoFileData) (models.PluginDesc, error) {
	var err error
	protoPlugin := models.PluginDesc{
		ModuleName: config.PluginModuleName,
	}

	processedServices := make(map[string]struct{})

	for _, protoFile := range protoFiles {
		logger.Infof(ctx, "Fetching %s file...", protoFile.OriginalPath)
		protoFile.RawData, err = c.fileFetcher.FetchFile(ctx, protoFile.OriginalPath)
		if err != nil {
			return models.PluginDesc{}, err
		}

		protoDir := filepath.Join(pluginDir, protoFile.PkgName)
		if err := os.MkdirAll(protoDir, fs.ModePerm); err != nil {
			return models.PluginDesc{}, err
		}

		if err := ioutil.WriteFile(filepath.Join(pluginDir, protoFile.Path), protoFile.RawData, fs.ModePerm); err != nil {
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
