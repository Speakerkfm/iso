package command

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

const (
	defaultPath = "."

	projectDir     = "proto"
	configFileName = "config.yaml"
)

func (c *Command) Init(path string) {
	configData, err := c.gen.GenerateConfig()
	if err != nil {
		handleError(err)
		return
	}

	if path == "" {
		path = defaultPath
	}

	if err := os.Mkdir(fmt.Sprintf("%s/%s", path, projectDir), fs.ModePerm); err != nil {
		handleError(err)
		return
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", path, projectDir, configFileName), configData, fs.ModePerm); err != nil {
		handleError(err)
		return
	}

	fmt.Fprintln(os.Stdout, "Project initialized")
}
