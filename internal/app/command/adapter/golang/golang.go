package golang

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

type Golang struct {
}

func New() *Golang {
	return &Golang{}
}

func (g *Golang) BuildPlugin(wd, outDir, moduleName, buildFile string) error {
	cmdInit := exec.Command("go", "mod", "init", moduleName)
	cmdInit.Dir = wd

	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("fail to init module: %w", err)
	}

	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = wd

	if err := cmdTidy.Run(); err != nil {
		return fmt.Errorf("fail to load deps: %w", err)
	}

	cmdBuildModule := exec.Command("go", "build", "-buildmode=plugin", "-o", outDir, buildFile)
	cmdBuildModule.Dir = wd

	logger.Infof(context.Background(), "Exec: %s", cmdBuildModule.String())
	if err := cmdBuildModule.Run(); err != nil {
		return fmt.Errorf("fail to build plugin: %w", err)
	}

	return nil
}
