package protoc

import (
	"os/exec"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Protoc struct {
}

func New() *Protoc {
	return &Protoc{}
}

func (p *Protoc) Process(protoFile *models.ProtoFile) error {
	cmd := exec.Command("protoc", "--go_out=.", "--go_opt=paths=source_relative", protoFile.Path)

	return cmd.Run()
}
