package models

import (
	"time"
)

// ExternalDependency сущность, которая хранит описание внешней зависимости
type ExternalDependency struct {
	Host       string   `yaml:"host"`
	ProtoPaths []string `yaml:"grpc,flow"`
}

// ServiceSpecification сущность, которая содержит описание характеристики изолируемого сервиса
type ServiceSpecification struct {
	ExternalDependencies []ExternalDependency `yaml:"external_dependencies,flow"`
}

type ServiceConfigDesc struct {
	Host         string              `yaml:"host"`
	GRPCHandlers []HandlerConfigDesc `yaml:"-"`
}

type HandlerConfigDesc struct {
	ServiceName string     `yaml:"service_name"`
	MethodName  string     `yaml:"method_name"`
	Rules       []RuleDesc `yaml:"rules,flow"`
}

type RuleDesc struct {
	Conditions []HandlerConditionDesc `yaml:"conditions,flow"`
	Response   HandlerResponseDesc    `yaml:"response"`
}

type HandlerResponseDesc struct {
	Delay time.Duration `yaml:"delay"`
	Data  string        `yaml:"data"`
}

type HandlerConditionDesc struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
