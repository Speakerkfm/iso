package generator

import (
	_ "embed"
	"text/template"
)

//go:embed source/spec_example.yaml
var specExampleSource []byte

//go:embed source/spec_plugin.tmpl
var specPluginTemplateSource string

//go:embed source/reverse_proxy.tmpl
var reverseProxyConfigTemplateSource string

var specPluginTemplate = template.Must(template.New("spec_plugin").Parse(specPluginTemplateSource))

var reverseProxyConfigTemplate = template.Must(template.New("reverse_proxy").Parse(reverseProxyConfigTemplateSource))
