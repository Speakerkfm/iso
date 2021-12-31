package command

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

const (
	defaultPath = "."

	configFileName = "config.yaml"
)

func (c *Command) Init(ctx context.Context, path string) {
	configData, err := c.gen.GenerateConfigData()
	if err != nil {
		handleError(err)
		return
	}

	if path == "" {
		path = defaultPath
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", path, configFileName), configData, fs.ModePerm); err != nil {
		handleError(err)
		return
	}

	fmt.Fprintln(os.Stdout, "Project initialized")
}
