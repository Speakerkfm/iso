package command

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Speakerkfm/iso/internal/pkg/config"
)

func (c *Command) RulesApply(ctx context.Context, dir string) error {
	rulesDir := filepath.Join(dir, config.RulesDir)
	svcConfigs, err := c.ruleParser.ParseDirectory(ctx, rulesDir)
	if err != nil {
		return fmt.Errorf("fail to parse service configs from direcory:%s, err: %w", rulesDir, err)
	}

	if err := c.isoSrv.SaveServiceConfigs(ctx, svcConfigs); err != nil {
		return fmt.Errorf("fail to save service configs: %s, err: %w", dir, err)
	}
	return nil
}
