package manager

import (
	"context"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/logger"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type valueGetter func(ctx context.Context, key string) (string, bool)

type manager struct {
	ruleTree *models.RuleNode
}

func New() *manager {
	return &manager{}
}

func (m *manager) GetRule(ctx context.Context, req models.Request) (*models.Rule, error) {
	logger.Infof(ctx, "Got request: %+v", req)

	currentNode := m.ruleTree
	for currentNode.Rule == nil {
		nextNodeFound := false
		for _, nextNode := range currentNode.NextNodes {
			if evalCondition(ctx, req.GetValue, nextNode.Condition) {
				nextNodeFound = true
				currentNode = nextNode
				break
			}
		}
		if !nextNodeFound {
			break
		}
	}
	if currentNode == nil || currentNode.Rule == nil {
		return nil, fmt.Errorf("rule not found")
	}
	return currentNode.Rule, nil
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

		leafFound := false
		for !leafFound {
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
				leafFound = true
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

func evalCondition(ctx context.Context, getValue valueGetter, cond models.Condition) bool {
	v, ok := getValue(ctx, cond.Key)
	if !ok {
		return false
	}
	if cond.Value == "*" {
		return true
	}
	return v == cond.Value
}
