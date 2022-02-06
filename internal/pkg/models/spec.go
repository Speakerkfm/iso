package models

// ExternalDependency сущность, которая хранит описание внешней зависимости из файла конфигурации
type ExternalDependency struct {
	Host       string   `yaml:"host"`
	ProtoPaths []string `yaml:"grpc,flow"`
}

// Config сущность, которая содержит описание из файла конфигурации
type Config struct {
	ExternalDependencies []ExternalDependency `yaml:"external_dependencies,flow"`
}
