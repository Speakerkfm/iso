package cobra

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/Speakerkfm/iso/internal/app/command"
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
		Short: "Iso is a tool for grpc mocking",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := c.Root(ctx); err != nil {
				handleError(err)
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
				handleError(err)
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
				handleError(err)
			}
		},
	}
}

func handleError(err error) {
	log.Fatal(err.Error())
}
