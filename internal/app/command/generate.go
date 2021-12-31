package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

const (
	protoPluginName = "struct.go"
)

func (c *Command) Generate(path string) {
	ctx := context.Background()
	workDir := fmt.Sprintf("%s/%s/pb", path, projectDir)

	spec, err := c.loadConfig(path)
	if err != nil {
		handleError(err)
	}

	protoFiles, err := c.processConfig(spec, workDir)
	if err != nil {
		handleError(err)
	}

	if err := os.MkdirAll(workDir, fs.ModePerm); err != nil {
		handleError(err)
	}

	svcDesc, err := c.processProtoFiles(ctx, protoFiles, workDir)
	if err != nil {
		handleError(err)
	}

	protoPluginData, err := c.gen.GenerateProtoPlugin(svcDesc)
	if err != nil {
		handleError(err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", path, projectDir, protoPluginName), protoPluginData, fs.ModePerm); err != nil {
		handleError(err)
		return
	}

	fmt.Fprintln(os.Stdout, "Project generated")
}

func (c *Command) loadConfig(path string) (models.Config, error) {
	spec := models.Config{}

	fin, err := os.Open(fmt.Sprintf("%s/%s/%s", path, projectDir, configFileName))
	if err != nil {
		return models.Config{}, err
	}
	defer fin.Close()

	if err := yaml.NewDecoder(fin).Decode(&spec); err != nil {
		return models.Config{}, err
	}

	return spec, nil
}

func (c *Command) processConfig(spec models.Config, workDir string) ([]*models.ProtoFile, error) {
	var protoFiles []*models.ProtoFile

	for _, dep := range spec.ExternalDependencies {
		for _, protoPath := range dep.ProtoPaths {
			fileName := filepath.Base(protoPath)
			pkgName := strings.Split(fileName, ".")[0]
			protoFiles = append(protoFiles, &models.ProtoFile{
				Name:         fileName,
				PkgName:      pkgName,
				OriginalPath: protoPath,
				Path:         fmt.Sprintf("%s/%s/%s", workDir, pkgName, fileName),
			})
		}
	}

	return protoFiles, nil
}

func (c *Command) processProtoFiles(ctx context.Context, protoFiles []*models.ProtoFile, workDir string) (map[string]*models.ProtoServiceDesc, error) {
	res := make(map[string]*models.ProtoServiceDesc)

	for _, protoFile := range protoFiles {
		protoData, err := c.fileFetcher.FetchFile(ctx, protoFile.OriginalPath)
		if err != nil {
			return nil, err
		}
		protoFile.RawData = protoData

		protoDir := fmt.Sprintf("%s/%s", workDir, protoFile.PkgName)

		if err := os.MkdirAll(protoDir, fs.ModePerm); err != nil {
			handleError(err)
		}

		if err := ioutil.WriteFile(protoFile.Path, protoData, fs.ModePerm); err != nil {
			return nil, err
		}

		if err := c.protoc.Process(protoFile); err != nil {
			return nil, err
		}

		serviceDescriptions, err := c.protoParser.Parse(protoFile.RawData)
		if err != nil {
			return nil, err
		}

		for _, svcDesc := range serviceDescriptions {
			svcDesc.ProtoPath = protoFile.Path
			res[svcDesc.Name] = svcDesc
		}
	}

	return res, nil
}
