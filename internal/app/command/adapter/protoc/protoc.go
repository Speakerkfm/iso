package protoc

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Protoc struct {
}

func New() *Protoc {
	return &Protoc{}
}

func (p *Protoc) Process(wd string, protoFile *models.ProtoFileData) error {
	cmd := exec.Command("protoc", "--go_out=.", "--go_opt=paths=source_relative", protoFile.Path)
	cmd.Dir = wd

	logger.Infof(context.Background(), "Exec: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fail to process protoc for file %s: %w", protoFile.Name, err)
	}

	return nil
}
