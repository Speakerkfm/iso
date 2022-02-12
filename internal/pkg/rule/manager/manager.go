package manager

import (
	"context"
	"encoding/json"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type manager struct {
	ruleTree *RuleNode
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
			if nextNode.Condition.Eval(req.Values) {
				currentNode = nextNode
				break
			}
		}
	}
	return currentNode.Rule, nil
}

type RuleNode struct {
	Condition models.Condition
	NextNodes []*RuleNode
	Rule      *models.Rule
}

func createRuleTree() *RuleNode {
	return &RuleNode{
		Condition: models.Condition{},
		NextNodes: []*RuleNode{
			{
				Condition: models.Condition{
					Key:   "Host",
					Value: "127.0.0.1",
				},
				NextNodes: []*RuleNode{
					{
						Condition: models.Condition{
							Key:   "ServiceName",
							Value: "UserService",
						},
						NextNodes: []*RuleNode{
							{
								Condition: models.Condition{
									Key:   "MethodName",
									Value: "GetUser",
								},
								NextNodes: []*RuleNode{
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
							Key:   "ServiceName",
							Value: "PhoneService",
						},
						NextNodes: []*RuleNode{
							{
								Condition: models.Condition{
									Key:   "MethodName",
									Value: "CheckPhone",
								},
								NextNodes: []*RuleNode{
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
			Key:   "Host",
			Value: "127.0.0.1",
		},
		{
			Key:   "ServiceName",
			Value: "UserService",
		},
		{
			Key:   "MethodName",
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
			Key:   "Host",
			Value: "127.0.0.1",
		},
		{
			Key:   "ServiceName",
			Value: "UserService",
		},
		{
			Key:   "MethodName",
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
			Key:   "Host",
			Value: "127.0.0.1",
		},
		{
			Key:   "ServiceName",
			Value: "PhoneService",
		},
		{
			Key:   "MethodName",
			Value: "CheckPhone",
		},
		{
			Key:   "body.phone",
			Value: "+1000",
		},
	},
	MessageData: json.RawMessage(`{"exists":true}`),
}
