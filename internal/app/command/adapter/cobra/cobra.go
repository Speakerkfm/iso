package cobra

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Speakerkfm/iso/internal/app/command"
	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

const (
	defaultPath = "."
)

func New(c *command.Command) *cobra.Command {
	root := handleRoot(c)
	init := handleInit(c)
	generate := handleGenerate(c)

	server := handleServer(c)
	serverStart := handleServerStart(c)
	serverStop := handleServerStop(c)
	server.AddCommand(serverStart)
	server.AddCommand(serverStop)

	rules := handleRules(c)
	rulesSync := handleRulesSync(c)
	rulesApply := handleRulesApply(c)
	rules.AddCommand(rulesSync)
	rules.AddCommand(rulesApply)

	report := handleReport(c)
	reportLoad := handleReportLoad(c)
	report.AddCommand(reportLoad)

	root.AddCommand(init)
	root.AddCommand(generate)
	root.AddCommand(server)
	root.AddCommand(rules)
	root.AddCommand(report)

	return root
}

func handleRoot(c *command.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "iso",
		Short: "Iso is a tool for mocking web interfaces by generated data from specification files",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := c.Root(ctx); err != nil {
				handleError(ctx, err)
			}
		},
	}
}

func handleInit(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init project",
		Long:  `Init project with specification file.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := config.DefaultProjectDir
			if len(args) > 0 {
				path = args[0]
			}

			ctx := context.Background()
			if err := c.Init(ctx, path); err != nil {
				handleError(ctx, err)
			}
		},
	}

	return cmd
}

func handleGenerate(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate project",
		Long:  `Generate project from specification file.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := config.DefaultProjectDir
			if len(args) > 0 {
				path = args[0]
			}
			dockerEnabled, _ := cmd.Flags().GetBool("docker")

			ctx := context.Background()
			if err := c.Generate(ctx, path, dockerEnabled); err != nil {
				handleError(ctx, err)
			}
		},
	}

	cmd.Flags().Bool("docker", false, "")

	return cmd
}

func handleServer(c *command.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  `Server commands.`,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}

func handleServerStart(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start server",
		Long:  `Start server background.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := config.DefaultProjectDir
			if len(args) > 0 {
				path = args[0]
			}
			dockerEnabled, _ := cmd.Flags().GetBool("docker")

			ctx := context.Background()
			if err := c.StartServer(ctx, path, dockerEnabled); err != nil {
				handleError(ctx, err)
			}
		},
	}

	cmd.Flags().Bool("docker", false, "")

	return cmd
}

func handleServerStop(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop server",
		Long:  `Stop server background.`,
		Run: func(cmd *cobra.Command, args []string) {
			dockerEnabled, _ := cmd.Flags().GetBool("docker")

			ctx := context.Background()
			if err := c.StopServer(ctx, dockerEnabled); err != nil {
				handleError(ctx, err)
			}
		},
	}

	cmd.Flags().Bool("docker", false, "")

	return cmd
}

func handleRules(c *command.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: "rules",
		Long:  `Rules commands.`,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}

func handleRulesSync(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync rules",
		Long:  `Sync rules with iso server.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := config.DefaultProjectDir
			if len(args) > 0 {
				path = args[0]
			}

			ctx := context.Background()
			if err := c.RulesSync(ctx, path); err != nil {
				handleError(ctx, err)
			}
		},
	}

	return cmd
}

func handleRulesApply(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply rules",
		Long:  `Apply rules to iso server.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := config.DefaultProjectDir
			if len(args) > 0 {
				path = args[0]
			}

			ctx := context.Background()
			if err := c.RulesApply(ctx, path); err != nil {
				handleError(ctx, err)
			}
		},
	}

	return cmd
}

func handleReport(c *command.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "report",
		Long:  `Report commands.`,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}

func handleReportLoad(c *command.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "lead report",
		Long:  `Load report from iso server.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := c.ReportLoad(ctx); err != nil {
				handleError(ctx, err)
			}
		},
	}

	cmd.Flags().Bool("docker", false, "")

	return cmd
}

func handleError(ctx context.Context, err error) {
	logger.Fatal(ctx, err.Error())
}
