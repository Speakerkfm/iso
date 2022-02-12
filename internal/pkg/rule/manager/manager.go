package manager

import (
	"context"
	"encoding/json"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type manager struct {
	ruleTree *models.RuleNode
}

func New() *manager {
	return &manager{
		ruleTree: createRuleTree(),
	}
}

func (m *manager) GetRule(ctx context.Context, req *models.Request) (*models.Rule, error) {
	currentNode := m.ruleTree
	for currentNode.Rule == nil {
		for _, nextNode := range currentNode.NextNodes {
			if evalCondition(nextNode.Condition, req.Values) {
				currentNode = nextNode
				break
			}
		}
	}
	return currentNode.Rule, nil
}

func evalCondition(cond models.Condition, values map[string]string) bool {
	v, ok := values[cond.Key]
	if !ok {
		return false
	}
	return v == cond.Value
}

func createRuleTree() *models.RuleNode {
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
											Value: "15",
										},
										Rule: rule2,
									},
									{
										Condition: models.Condition{
											Key:   "body.id",
											Value: "10",
										},
										Rule: rule1,
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
	MessageData: json.RawMessage(`{"user":{"id":10,"name":"kek_10"}}`),
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
	MessageData: json.RawMessage(`{"user":{"id":15,"name":"kek_15"}}`),
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
	MessageData: json.RawMessage(`{"exists":true}`),
}
