package command

import (
	"context"
	"fmt"
	"os"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Generator interface {
	GenerateConfig() ([]byte, error)
}

type FileFetcher interface {
	FetchFile(ctx context.Context, filePath string) ([]byte, error)
}

type Protoc interface {
	Process(protoFile *models.ProtoFile) error
}

type ProtoParser interface {
	Parse(rawProtoFile []byte) ([]*models.ProtoServiceDesc, error)
}

type Command struct {
	gen         Generator
	fileFetcher FileFetcher
	protoc      Protoc
	protoParser ProtoParser
}

func New(g Generator, ff FileFetcher, pc Protoc, pp ProtoParser) *Command {
	return &Command{
		gen:         g,
		fileFetcher: ff,
		protoc:      pc,
		protoParser: pp,
	}
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
