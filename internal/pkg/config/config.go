package config

import (
	"time"
)

const (
	ISOServerAdminHost   = "0.0.0.0:8150"
	ISOServerGRPCHost    = "localhost:8152"
	ISOServerDockerImage = "speakerkfm/iso-server:latest"
	ISOServerDockerID    = "my-iso-server"

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
	PluginDockerImage = "speakerkfm/iso-plugin:latest"

	DefaultProjectDir     = "."
	SpecificationFileName = "spec.yaml"

	RulesDir                   = "rules"
	ServiceConfigFileName      = "service.yaml"
	ProtoHandlerDirName        = "proto"
	ReverseProxyConfigFileName = "iso_nginx.conf"

	HandlerConfigDefaultTimeout = 5 * time.Millisecond
	RulesSyncInterval           = 5 * time.Second

	BatcherEnabled          = true
	BatcherBatchCount       = 3
	BatcherFlushInterval    = 1 * time.Second
	BatcherFlushItemsAmount = 100
	BatcherEventBuffSize    = 1000
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
