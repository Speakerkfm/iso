package manager

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

func Test_createRuleTree(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expected := createTestRuleTree()
		actual := createRuleTree([]*models.Rule{rule1, rule2, rule3})

		assert.Equal(t, expected, actual)
	})
}

func TestManager_GetHandlerConfig(t *testing.T) {
	t.Run("ok", func(t *testing.T) {

	})
}

var rule1 = &models.Rule{
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
}

var rule2 = &models.Rule{
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
}

var rule3 = &models.Rule{
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
}

func createTestRuleTree() *models.RuleNode {
	return &models.RuleNode{
		Condition: models.Condition{},
		NextNodes: []*models.RuleNode{
			{
				Condition: models.Condition{
					Key:   models.FieldHost,
					Value: "127.0.0.1",
				},
				NextNodes: []*models.RuleNode{
					{
						Condition: models.Condition{
							Key:   models.FieldServiceName,
							Value: "UserService",
						},
						NextNodes: []*models.RuleNode{
							{
								Condition: models.Condition{
									Key:   models.FieldMethodName,
									Value: "GetUser",
								},
								NextNodes: []*models.RuleNode{
									{
										Condition: models.Condition{
											Key:   "body.id",
											Value: "10",
										},
										Rule: rule1,
									},
									{
										Condition: models.Condition{
											Key:   "body.id",
											Value: "15",
										},
										Rule: rule2,
									},
								},
							},
						},
					},
					{
						Condition: models.Condition{
							Key:   models.FieldServiceName,
							Value: "PhoneService",
						},
						NextNodes: []*models.RuleNode{
							{
								Condition: models.Condition{
									Key:   models.FieldMethodName,
									Value: "CheckPhone",
								},
								NextNodes: []*models.RuleNode{
									{
										Condition: models.Condition{
											Key:   "body.phone",
											Value: "+1000",
										},
										Rule: rule3,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
