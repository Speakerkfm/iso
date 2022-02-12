package models

import (
	"encoding/json"
)

type Rule struct {
	Conditions  []Condition
	MessageData json.RawMessage
}

type Condition struct {
	Key   string
	Value string
}

func (c Condition) Eval(values map[string]string) bool {
	v, ok := values[c.Key]
	if !ok {
		return false
	}
	return v == c.Value
}
