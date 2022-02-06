package generator

import (
	_ "embed"
	"text/template"
)

//go:embed source/config_example.yaml
var configTemplateExample []byte

//go:embed source/impl.tmpl
var implTemplateSource string

var implTemplate = template.Must(template.New("impl").Parse(implTemplateSource))
