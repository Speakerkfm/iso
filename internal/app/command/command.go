package command

import (
	"github.com/Speakerkfm/iso/internal/pkg/fetcher"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/parser/proto"
)

type Protoc interface {
	Process(wd string, protoFile *models.ProtoFileData) error
}

type Golang interface {
	BuildPlugin(wd, outDir, modName, buildFile string) error
}

type Command struct {
	gen         generator.Generator
	fileFetcher fetcher.FileFetcher
	protoParser proto.Parser
	protoc      Protoc
	golang      Golang
}

func New(g generator.Generator, ff fetcher.FileFetcher, pc Protoc, pp proto.Parser, golang Golang) *Command {
	return &Command{
		gen:         g,
		fileFetcher: ff,
		protoc:      pc,
		protoParser: pp,
		golang:      golang,
	}
}
