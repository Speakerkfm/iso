package command

import (
	"context"

	"github.com/Speakerkfm/iso/internal/pkg/fetcher"
	"github.com/Speakerkfm/iso/internal/pkg/generator"
	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/parser/proto"
	rule_parser "github.com/Speakerkfm/iso/internal/pkg/rule/parser"
)

type Protoc interface {
	Process(wd string, protoFile *models.ProtoFileData) error
}

type Golang interface {
	BuildPlugin(wd, outDir, modName, buildFile string) error
}

type Docker interface {
	StartServer(dir string) error
	StopServer() error
	BuildPlugin(wd, outDir, modName, buildFile string) error
}

type ISOServer interface {
	GetServiceConfigs(ctx context.Context) ([]models.ServiceConfigDesc, error)
	GetReport(ctx context.Context) (*models.Report, error)
	SaveServiceConfigs(ctx context.Context, serviceConfigs []models.ServiceConfigDesc) error
}

type Command struct {
	gen         generator.Generator
	fileFetcher fetcher.FileFetcher
	protoParser proto.Parser
	ruleParser  rule_parser.Parser
	protoc      Protoc
	golang      Golang
	docker      Docker
	isoSrv      ISOServer
}

func New(gen generator.Generator,
	fileFetcher fetcher.FileFetcher,
	protoc Protoc,
	protoParser proto.Parser,
	ruleParser rule_parser.Parser,
	golang Golang,
	docker Docker,
	isoSrv ISOServer) *Command {
	return &Command{
		gen:         gen,
		fileFetcher: fileFetcher,
		protoc:      protoc,
		protoParser: protoParser,
		ruleParser:  ruleParser,
		golang:      golang,
		docker:      docker,
		isoSrv:      isoSrv,
	}
}
