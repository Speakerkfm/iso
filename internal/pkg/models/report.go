package models

type MethodReport struct {
	RuleStat map[string]int64 `json:"rule_stat"`
}

type ServiceReport struct {
	Method map[string]*MethodReport `json:"method"`
}

type Report struct {
	Service map[string]*ServiceReport `json:"service"`
}
