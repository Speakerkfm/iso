package models

type MethodStat struct {
	SuccessCount int `json:"success_count"`
	ErrorCount   int `json:"error_count"`
}

type MethodReport struct {
	Stat *MethodStat `json:"stat"`
}

type ServiceReport struct {
	Method map[string]*MethodReport `json:"method"`
}

type Report struct {
	Service map[string]*ServiceReport `json:"service"`
}
