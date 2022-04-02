package admin

import (
	"context"
	"net/http"

	"github.com/Speakerkfm/iso/internal/pkg/router"
	"github.com/Speakerkfm/iso/internal/pkg/rule/service"
)

type Implementation struct {
	ruleSvc service.RuleService
}

func New(ruleSvc service.RuleService) *Implementation {
	return &Implementation{
		ruleSvc: ruleSvc,
	}
}

func (i *Implementation) RegisterGateway(ctx context.Context, mux router.ServeMux) error {
	mux.MethodFunc(http.MethodGet, "/service_configs", i.HandleGetServiceConfigs)
	mux.MethodFunc(http.MethodPut, "/service_configs", i.HandleSaveServiceConfigs)
	return nil
}
