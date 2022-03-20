package main

import (
	"context"

	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"

	"github.com/Speakerkfm/iso/internal/app/command"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/cobra"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/golang"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/protoc"
	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/fetcher"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/parser/proto"
)

const (
	configPath = "" // now config with constants
)

func main() {
	appCtx := context.Background()
	if err := config.Parse(configPath); err != nil {
		logger.Fatalf(appCtx, "fail to parse config file: %s", err.Error())
	}

	g := generator.New()
	ff := fetcher.New()
	pc := protoc.New()
	pp := proto.New()
	glng := golang.New()

	cmd := command.New(g, ff, pc, pp, glng)

	cobraCmd := cobra.New(cmd)

	if err := cobraCmd.Execute(); err != nil {
		logger.Fatalf(appCtx, "fail to execute command: %s", err.Error())
	}
}
