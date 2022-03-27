package config

import (
	"time"
)

const (
	ISOServerHost        = "localhost:8150"
	ISOServerDockerImage = "iso-server"

	RequestHeaderHost  = "x-original-host"
	RequestHeaderReqID = "x-request-id"

	RequestFieldHost        = "header.x-original-host"
	RequestFieldReqID       = "header.x-request-id"
	RequestFieldServiceName = "ServiceName"
	RequestFieldMethodName  = "MethodName"

	PluginFileName    = "spec.so"
	PluginGoFileName  = "spec.go"
	PluginModuleName  = "iso_plugin"
	PluginDir         = "plugin"
	PluginDockerImage = "iso-plugin"

	DefaultProjectDir     = "."
	SpecificationFileName = "spec.yaml"

	RulesDir                   = "rules"
	ServiceConfigFileName      = "service.yaml"
	ProtoHandlerDirName        = "proto"
	ReverseProxyConfigFileName = "iso_nginx.conf"

	HandlerConfigDefaultTimeout = 5 * time.Millisecond
)

// Parse парсит конфигурационный файл и заполняет структуры конфига
func Parse(cfgPath string) error {
	// fin, err := os.Open(cfgPath)
	// if err != nil {
	// 	return err
	// }
	// defer fin.Close()
	//
	// if err := yaml.NewDecoder(fin).Decode(&cfg); err != nil {
	// 	return err
	// }

	return nil
}
