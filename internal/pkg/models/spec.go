package models

import (
	"time"
)

// PluginDesc сущность, которая хранит объекты для генерации плагина
type PluginDesc struct {
	ModuleName    string
	Imports       []string
	ProtoServices []*ProtoServiceDesc
}

// ExternalDependency сущность, которая хранит описание внешней зависимости
type ExternalDependency struct {
	Host       string   `yaml:"host"`
	Name       string   `yaml:"name"`
	ProtoPaths []string `yaml:"proto,flow"`
}

// ServiceSpecification сущность, которая содержит описание характеристики изолируемого сервиса
type ServiceSpecification struct {
	ExternalDependencies []ExternalDependency `yaml:"external_dependencies,flow"`
}

type ServiceConfigDesc struct {
	Host         string              `yaml:"host"`
	Name         string              `yaml:"-"`
	GRPCHandlers []HandlerConfigDesc `yaml:"-"`
}

type HandlerConfigDesc struct {
	ServiceName string     `yaml:"service_name"`
	MethodName  string     `yaml:"method_name"`
	Rules       []RuleDesc `yaml:"rules"`
}

type RuleDesc struct {
	Conditions []HandlerConditionDesc `yaml:"conditions"`
	Response   HandlerResponseDesc    `yaml:"response"`
}

type HandlerResponseDesc struct {
	Delay time.Duration `yaml:"delay"`
	Data  string        `yaml:"data"`
	Error string        `yaml:"error"`
}

type HandlerConditionDesc struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
