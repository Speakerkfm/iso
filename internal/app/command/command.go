package command

import (
	"github.com/Speakerkfm/iso/internal/pkg/fetcher"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/proto_parser"
)

type Protoc interface {
	Process(wd string, protoFile *models.ProtoFileData) error
}

type Golang interface {
	CreateModule(wd, modName string) error
	BuildPlugin(wd, outDir, buildFile string) error
}

type Command struct {
	gen         generator.Generator
	fileFetcher fetcher.FileFetcher
	protoParser proto_parser.Parser
	protoc      Protoc
	golang      Golang
}

func New(g generator.Generator, ff fetcher.FileFetcher, pc Protoc, pp proto_parser.Parser, golang Golang) *Command {
	return &Command{
		gen:         g,
		fileFetcher: ff,
		protoc:      pc,
		protoParser: pp,
		golang:      golang,
	}
}
