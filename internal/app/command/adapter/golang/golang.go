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

func (g *Golang) CreateModule(wd, modName string) error {
	cmdInit := exec.Command("go", "mod", "init", modName)
	cmdInit.Dir = wd

	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("fail to init module: %w", err)
	}

	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = wd

	if err := cmdTidy.Run(); err != nil {
		return fmt.Errorf("fail to load deps: %w", err)
	}

	return nil
}

func (g *Golang) BuildPlugin(wd, buildFile string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to get working dir: %w", err)
	}

	cmdBuildModule := exec.Command("go", "build", "-buildmode=plugin", "-o", currentDir, buildFile)
	cmdBuildModule.Dir = wd

	if err := cmdBuildModule.Run(); err != nil {
		return fmt.Errorf("fail to build plugin: %w", err)
	}

	return nil
}
