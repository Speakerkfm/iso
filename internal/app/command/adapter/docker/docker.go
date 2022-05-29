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
		"--env", fmt.Sprintf("PLUGIN_FILE_NAME=%s", config.PluginFileName),
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
	cmdStartServer := exec.Command("docker", "run",
		"-d",
		"--name", config.ISOServerDockerID,
		"-v", fmt.Sprintf("%s:/iso", dir),
		"-p", "82:82",
		"-p", "8150:8150",
		config.ISOServerDockerImage)
	cmdStartServer.Dir = dir

	logger.Infof(context.Background(), "Exec: %s", cmdStartServer.String())
	if err := cmdStartServer.Run(); err != nil {
		return fmt.Errorf("fail to start server: %w", err)
	}

	return nil
}

// StartServer ...
func (d *Docker) StopServer() error {
	cmdStartServer := exec.Command("docker", "rm", "-f", config.ISOServerDockerID)

	logger.Infof(context.Background(), "Exec: %s", cmdStartServer.String())
	if err := cmdStartServer.Run(); err != nil {
		return fmt.Errorf("fail to stop server: %w", err)
	}

	return nil
}
