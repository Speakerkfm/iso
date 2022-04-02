package generator

import (
	_ "embed"
	"text/template"
)

//go:embed source/spec_example.yaml
var specExampleSource []byte

//go:embed source/spec_plugin.tmpl
var specPluginTemplateSource string

var specPluginTemplate = template.Must(template.New("spec_plugin").Parse(specPluginTemplateSource))
