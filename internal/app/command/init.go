package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
)

const (
	defaultPath = "."

	configFileName = "config.yaml"
)

func (c *Command) Init(ctx context.Context, path string) error {
	configData, err := c.gen.GenerateConfigData()
	if err != nil {
		return fmt.Errorf("fail to generate config data: %w", err)
	}

	if path == "" {
		path = defaultPath
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", path, configFileName), configData, fs.ModePerm); err != nil {
		return fmt.Errorf("fail to save config data to file: %w", err)
	}

	log.Println( "Project initialized")

	return nil
}
