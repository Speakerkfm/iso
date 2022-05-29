package models

import (
	"encoding/json"
	"time"
)

// RuleNode узел дерева правил
type RuleNode struct {
	Condition Condition
	NextNodes []*RuleNode
	Rule      *Rule
}

// Rule правило для ответа
type Rule struct {
	ID            string
	Name          string
	Conditions    []Condition
	HandlerConfig *HandlerConfig
}

// HandlerConfig конфигурация обработчика запроса
type HandlerConfig struct {
	ServiceName   string
	MethodName    string
	ResponseDelay time.Duration
	MessageData   json.RawMessage
	Error         string
}

// Condition условие правила
type Condition struct {
	Key   string
	Value string
}
