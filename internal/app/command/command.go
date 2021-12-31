package command

import (
	"context"
	"fmt"
	"os"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Generator interface {
	GenerateConfigData() ([]byte, error)
	GenerateProtoPluginData(protoPlugin models.ProtoPlugin) ([]byte, error)
}

type FileFetcher interface {
	FetchFile(ctx context.Context, filePath string) ([]byte, error)
}

type Protoc interface {
	Process(wd string, protoFile *models.ProtoFile) error
}

type ProtoParser interface {
	Parse(rawProtoFile []byte) ([]*models.ProtoServiceDesc, error)
}

type Golang interface {
	CreateModule(wd, modName string) error
	BuildPlugin(wd, buildFile string) error
}

type Command struct {
	gen         Generator
	fileFetcher FileFetcher
	protoc      Protoc
	protoParser ProtoParser
	golang      Golang
}

func New(g Generator, ff FileFetcher, pc Protoc, pp ProtoParser, golang Golang) *Command {
	return &Command{
		gen:         g,
		fileFetcher: ff,
		protoc:      pc,
		protoParser: pp,
		golang:      golang,
	}
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
