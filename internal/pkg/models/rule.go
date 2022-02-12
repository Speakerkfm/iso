package models

import (
	"encoding/json"
)

type RuleNode struct {
	Condition Condition
	NextNodes []*RuleNode
	Rule      *Rule
}

type Rule struct {
	Conditions  []Condition
	MessageData json.RawMessage
}

type Condition struct {
	Key   string
	Value string
}
