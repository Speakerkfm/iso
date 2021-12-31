package golang

import (
	"fmt"
	"os"
	"os/exec"
)

type Golang struct {
}

func New() *Golang {
	return &Golang{}
}

func (g *Golang) CreateModule(modName string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to get working dir: %w", err)
	}
	cmdInit := exec.Command("go", "mod", "init", modName)
	cmdInit.Dir = fmt.Sprintf("%s/%s", workingDir, modName)

	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("fail to init module: %w", err)
	}

	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = fmt.Sprintf("%s/%s", workingDir, modName)

	if err := cmdTidy.Run(); err != nil {
		return fmt.Errorf("fail to load deps: %w", err)
	}

	return nil
}

func (g *Golang) BuildPlugin(path, buildFile string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to get working dir: %w", err)
	}
	cmdBuildModule := exec.Command("go", "build", "-buildmode=plugin", buildFile)
	cmdBuildModule.Dir = fmt.Sprintf("%s/%s", workingDir, path)

	if err := cmdBuildModule.Run(); err != nil {
		return fmt.Errorf("fail to build plugin: %w", err)
	}

	return nil
}
