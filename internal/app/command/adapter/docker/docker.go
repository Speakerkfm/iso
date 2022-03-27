package docker

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

type Docker struct {
}

func New() *Docker {
	return &Docker{}
}

// BuildPlugin ...
func (d *Docker) BuildPlugin(dir, outDir, moduleName, buildFile string) error {
	cmdBuildModule := exec.Command("docker", "run",
		"--rm",
		"-v", fmt.Sprintf("%s:/iso", outDir),
		"--env", fmt.Sprintf("PLUGIN_DIR=%s", config.PluginDir),
		"--env", fmt.Sprintf("PLUGIN_MODULE_NAME=%s", moduleName),
		"--env", fmt.Sprintf("PLUGIN_FILE_NAME=%s", config.PluginServerFileName),
		"--env", fmt.Sprintf("PLUGIN_GO_FILE_NAME=%s", buildFile),
		config.PluginDockerImage)
	cmdBuildModule.Dir = dir

	logger.Infof(context.Background(), "Exec: %s", cmdBuildModule.String())
	if err := cmdBuildModule.Run(); err != nil {
		return fmt.Errorf("fail to build plugin: %w", err)
	}

	return nil
}

// StartServer ...
func (d *Docker) StartServer(dir string) error {
	cmdBuildModule := exec.Command("docker", "run",
		"-d",
		"-v", fmt.Sprintf("%s:/iso", dir),
		"--env", fmt.Sprintf("REVERSE_PROXY_CONFIG_FILE_NAME=%s", config.ReverseProxyConfigFileName),
		"-p", "82:82",
		config.ISOServerDockerImage)
	cmdBuildModule.Dir = dir

	logger.Infof(context.Background(), "Exec: %s", cmdBuildModule.String())
	if err := cmdBuildModule.Run(); err != nil {
		return fmt.Errorf("fail to build plugin: %w", err)
	}

	return nil
}
