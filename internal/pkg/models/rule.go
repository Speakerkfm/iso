package models

import (
	"encoding/json"
	"time"
)

type RuleNode struct {
	Condition Condition
	NextNodes []*RuleNode
	Rule      *Rule
}

type Rule struct {
	Conditions    []Condition
	HandlerConfig *HandlerConfig
}

type HandlerConfig struct {
	ResponseDelay time.Duration
	MessageData   json.RawMessage
	Error         error
}

type Condition struct {
	Key   string
	Value string
}
