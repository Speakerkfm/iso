package manager

import (
	"context"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type manager struct {
	ruleTree *models.RuleNode
}

func New() *manager {
	return &manager{}
}

func (m *manager) GetHandlerConfig(ctx context.Context, req *models.Request) (*models.HandlerConfig, error) {
	currentNode := m.ruleTree
	for currentNode.Rule == nil {
		nextNodeFound := false
		for _, nextNode := range currentNode.NextNodes {
			if evalCondition(nextNode.Condition, req.Values) {
				nextNodeFound = true
				currentNode = nextNode
				break
			}
		}
		if !nextNodeFound {
			break
		}
	}
	if currentNode == nil || currentNode.Rule == nil || currentNode.Rule.HandlerConfig == nil {
		return nil, fmt.Errorf("rule not found")
	}
	return currentNode.Rule.HandlerConfig, nil
}

func (m *manager) UpdateRuleTree(rules []*models.Rule) {
	ruleTree := createRuleTree(rules)

	m.ruleTree = ruleTree
}

func createRuleTree(rules []*models.Rule) *models.RuleNode {
	headNode := &models.RuleNode{}
	for _, rule := range rules {
		branch := ruleToTreeBranch(rule)
		currentTreeNode := headNode
		currentBranchNode := branch

		leastFound := false
		for !leastFound {
			nextNodeFound := false
			for _, nextNode := range currentTreeNode.NextNodes {
				if isConditionEqual(nextNode.Condition, currentBranchNode.Condition) {
					nextNodeFound = true
					currentTreeNode = nextNode
					currentBranchNode = currentBranchNode.NextNodes[0]
					break
				}
			}
			if !nextNodeFound {
				leastFound = true
			}
		}

		currentTreeNode.NextNodes = append(currentTreeNode.NextNodes, currentBranchNode)
	}

	return headNode
}

func isConditionEqual(first, second models.Condition) bool {
	return first == second
}

func ruleToTreeBranch(rule *models.Rule) *models.RuleNode {
	var previousNode *models.RuleNode = nil
	for idx := len(rule.Conditions) - 1; idx >= 0; idx-- {
		node := &models.RuleNode{
			Condition: rule.Conditions[idx],
		}
		if previousNode == nil {
			node.Rule = rule
		} else {
			node.NextNodes = []*models.RuleNode{previousNode}
		}
		previousNode = node
	}
	return previousNode
}

func evalCondition(cond models.Condition, values map[string]string) bool {
	v, ok := values[cond.Key]
	if !ok {
		return false
	}
	return v == cond.Value
}
