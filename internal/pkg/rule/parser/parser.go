package parser

import (
	"context"
	"encoding/json"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type parser struct {
}

func New() *parser {
	return &parser{}
}

func (p *parser) Parse(ctx context.Context, directoryPath string) ([]*models.Rule, error) {
	// parse directory for rules

	return []*models.Rule{
		{
			Conditions: []models.Condition{
				{
					Key:   models.FieldHost,
					Value: "127.0.0.1",
				},
				{
					Key:   models.FieldServiceName,
					Value: "UserService",
				},
				{
					Key:   models.FieldMethodName,
					Value: "GetUser",
				},
				{
					Key:   "body.id",
					Value: "10",
				},
			},
			HandlerConfig: &models.HandlerConfig{
				MessageData: json.RawMessage(`{"user":{"id":10,"name":"kek_10"}}`),
			},
		},
		{
			Conditions: []models.Condition{
				{
					Key:   models.FieldHost,
					Value: "127.0.0.1",
				},
				{
					Key:   models.FieldServiceName,
					Value: "UserService",
				},
				{
					Key:   models.FieldMethodName,
					Value: "GetUser",
				},
				{
					Key:   "body.id",
					Value: "15",
				},
			},
			HandlerConfig: &models.HandlerConfig{
				MessageData: json.RawMessage(`{"user":{"id":15,"name":"kek_15"}}`),
			},
		},
		{
			Conditions: []models.Condition{
				{
					Key:   models.FieldHost,
					Value: "127.0.0.1",
				},
				{
					Key:   models.FieldServiceName,
					Value: "PhoneService",
				},
				{
					Key:   models.FieldMethodName,
					Value: "CheckPhone",
				},
				{
					Key:   "body.phone",
					Value: "+1000",
				},
			},
			HandlerConfig: &models.HandlerConfig{
				MessageData: json.RawMessage(`{"exists":true}`),
			},
		},
	}, nil
}
