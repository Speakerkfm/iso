package service

import (
	"context"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type RuleService interface {
	GetServiceConfigs(ctx context.Context) ([]models.ServiceConfigDesc, error)
	SaveServiceConfigs(ctx context.Context, svcConfigs []models.ServiceConfigDesc) error
}

type Repository interface {
	GetServiceConfigs(ctx context.Context) ([]models.ServiceConfigDesc, error)
	SaveServiceConfigs(ctx context.Context, svcConfigs []models.ServiceConfigDesc) error
}

type Service struct {
	repo Repository
	gen  Generator
}

type Generator interface {
	GenerateRules(svcConfigs []models.ServiceConfigDesc) []*models.Rule
}

func New(repo Repository, gen Generator) *Service {
	return &Service{
		repo: repo,
		gen:  gen,
	}
}

// GetServiceConfigs ...
func (svc *Service) GetServiceConfigs(ctx context.Context) ([]models.ServiceConfigDesc, error) {
	return svc.repo.GetServiceConfigs(ctx)
}

// GetServiceConfigs ...
func (svc *Service) SaveServiceConfigs(ctx context.Context, svcConfigs []models.ServiceConfigDesc) error {
	return svc.repo.SaveServiceConfigs(ctx, svcConfigs)
}

// GetRules ...
func (svc *Service) GetRules(ctx context.Context) ([]*models.Rule, error) {
	serviceConfigs, err := svc.repo.GetServiceConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get service configs: %w", err)
	}
	rules := svc.gen.GenerateRules(serviceConfigs)

	return rules, nil
}
