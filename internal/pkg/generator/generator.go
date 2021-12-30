package generator

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateConfig() ([]byte, error) {
	return []byte(configTemplate), nil
}
