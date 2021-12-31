package models

import (
	"google.golang.org/protobuf/proto"
)

type ProtoService struct {
	Name      string
	Methods   []ProtoMethod
	ProtoPath string
	PkgName   string
}

type ProtoMethod struct {
	Name           string
	RequestStruct  proto.Message
	ResponseStruct proto.Message
}
