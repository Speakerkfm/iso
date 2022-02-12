package models

import (
	"google.golang.org/protobuf/proto"
)

const (
	ServiceProviderName = "ServiceProvider"
)

type ServiceProvider interface {
	GetList() []*ProtoService
}

type ProtoService struct {
	Name      string
	Methods   []ProtoMethod
	ProtoPath string
}

type ProtoMethod struct {
	Name           string
	RequestStruct  proto.Message
	ResponseStruct proto.Message
}
