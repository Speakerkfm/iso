package cobra

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Speakerkfm/iso/internal/app/command"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

func New(c *command.Command) *cobra.Command {
	root := handleRoot(c)
	init := handleInit(c)
	generate := handleGenerate(c)

	root.AddCommand(init)
	root.AddCommand(generate)

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
	return &cobra.Command{
		Use:   "init",
		Short: "Init project",
		Long:  `Init project with specification file.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			ctx := context.Background()
			if err := c.Init(ctx, path); err != nil {
				handleError(ctx, err)
			}
		},
	}
}

func handleGenerate(c *command.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate project",
		Long:  `Generate project from specification file.`,
		Run: func(cmd *cobra.Command, args []string) {
			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			ctx := context.Background()
			if err := c.Generate(ctx, path); err != nil {
				handleError(ctx, err)
			}
		},
	}
}

func handleError(ctx context.Context, err error) {
	logger.Fatal(ctx, err.Error())
}
