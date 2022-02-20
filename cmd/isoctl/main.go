package main

import (
	"context"

	"github.com/Speakerkfm/iso/internal/app/command"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/cobra"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/golang"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/protoc"
	"github.com/Speakerkfm/iso/internal/pkg/fetcher"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/proto_parser"
)

func main() {
	appCtx := context.Background()
	g := generator.New()
	ff := fetcher.New()
	pc := protoc.New()
	pp := proto_parser.New()
	glng := golang.New()

	cmd := command.New(g, ff, pc, pp, glng)

	cobraCmd := cobra.New(cmd)

	if err := cobraCmd.Execute(); err != nil {
		logger.Fatalf(appCtx, "fail to execute command: %s", err.Error())
	}
}
