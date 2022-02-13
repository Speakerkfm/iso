package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/util"
)

const (
	protoPluginFile = "struct.go"
	protoPluginName = "iso_proto"
)

func (c *Command) Generate(ctx context.Context, configPath string) error {
	log.Println("Loading config...")
	spec, err := c.loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("fail to load config from path %s: %w", configPath, err)
	}

	protoFiles, err := c.processConfig(spec)
	if err != nil {
		return fmt.Errorf("fail to process config: %w", err)
	}

	wd := fmt.Sprintf("%s%s", os.TempDir(), util.NewUUID())

	log.Printf("Working directory: %s\n", wd)

	if err := os.MkdirAll(wd, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to make temp dir %s: %w", wd, err)
	}

	log.Println("Processing proto files...")

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

	log.Println("Building proto plugin...")
	if err := c.golang.BuildPlugin(wd, protoPluginFile); err != nil {
		return fmt.Errorf("fail to build proto plugin: %w", err)
	}

	log.Println("Proto plugin generated")

	return nil
}

func (c *Command) loadConfig(path string) (models.Config, error) {
	spec := models.Config{}

	fin, err := os.Open(path)
	if err != nil {
		return models.Config{}, err
	}
	defer fin.Close()

	if err := yaml.NewDecoder(fin).Decode(&spec); err != nil {
		return models.Config{}, err
	}

	return spec, nil
}

func (c *Command) processConfig(spec models.Config) ([]*models.ProtoFile, error) {
	var protoFiles []*models.ProtoFile

	for _, dep := range spec.ExternalDependencies {
		for _, protoPath := range dep.ProtoPaths {
			fileName := filepath.Base(protoPath)
			pkgName := strings.Split(fileName, ".")[0]
			protoFiles = append(protoFiles, &models.ProtoFile{
				Name:         fileName,
				PkgName:      pkgName,
				OriginalPath: protoPath,
				Path:         fmt.Sprintf("%s/%s", pkgName, fileName),
			})
		}
	}

	return protoFiles, nil
}

func (c *Command) processProtoFiles(ctx context.Context, wd string, protoFiles []*models.ProtoFile) (models.ProtoPlugin, error) {
	var err error
	protoPlugin := models.ProtoPlugin{
		ModuleName: protoPluginName,
	}

	processedServices := make(map[string]struct{})

	for _, protoFile := range protoFiles {
		protoFile.RawData, err = c.fileFetcher.FetchFile(ctx, protoFile.OriginalPath)
		if err != nil {
			return models.ProtoPlugin{}, err
		}

		protoDir := fmt.Sprintf("%s/%s", wd, protoFile.PkgName)
		if err := os.MkdirAll(protoDir, fs.ModePerm); err != nil {
			return models.ProtoPlugin{}, err
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", wd, protoFile.Path), protoFile.RawData, fs.ModePerm); err != nil {
			return models.ProtoPlugin{}, err
		}

		if err := c.protoc.Process(wd, protoFile); err != nil {
			return models.ProtoPlugin{}, err
		}

		serviceDescriptions, err := c.protoParser.Parse(protoFile.RawData)
		if err != nil {
			return models.ProtoPlugin{}, err
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
