package models

type ExternalDependency struct {
	Host       string   `yaml:"host"`
	ProtoPaths []string `yaml:"grpc,flow"`
}

type Config struct {
	ExternalDependencies []ExternalDependency `yaml:"external_dependencies,flow"`
}
