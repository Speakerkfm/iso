package grpc

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type codec struct {
}

func (codec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(*Response)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	if vv.lazyBytes != nil {
		return vv.lazyBytes, nil
	}
	return proto.Marshal(vv.msg)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(*Request)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv.Msg)
}

func (codec) Name() string {
	return ""
}
