package main

import (
	"log"

	"github.com/Speakerkfm/iso/internal/app/command"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/cobra"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/golang"
	"github.com/Speakerkfm/iso/internal/app/command/adapter/protoc"
	"github.com/Speakerkfm/iso/internal/pkg/fetcher"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/proto_parser"
)

func main() {
	g := generator.New()
	ff := fetcher.New()
	pc := protoc.New()
	pp := proto_parser.New()
	glng := golang.New()

	cmd := command.New(g, ff, pc, pp, glng)

	cobraCmd := cobra.New(cmd)

	if err := cobraCmd.Execute(); err != nil {
		log.Fatalf("fail to execute command: %s", err.Error())
	}
}
