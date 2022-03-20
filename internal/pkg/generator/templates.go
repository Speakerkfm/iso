package generator

import (
	_ "embed"
	"text/template"
)

//go:embed source/config_example.yaml
var configTemplateExample []byte

//go:embed source/impl.tmpl
var implTemplateSource string

//go:embed source/reverse_proxy.tmpl
var reverseProxyConfigTemplateSource string

var implTemplate = template.Must(template.New("impl").Parse(implTemplateSource))

var reverseProxyConfigTemplate = template.Must(template.New("reverse_proxy").Parse(reverseProxyConfigTemplateSource))
