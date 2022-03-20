package grpc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/logger"
)

const (
	valueHeaderPrefix = "header."
)

type Request struct {
	Context     context.Context
	ServiceName string
	MethodName  string
	Msg         proto.Message
	Headers     map[string][]string
	Values      map[string]string
}

func (r *Request) GetValue(ctx context.Context, key string) (string, bool) {
	if strings.HasPrefix(key, valueHeaderPrefix) {
		return r.getHeader(ctx, strings.TrimPrefix(key, valueHeaderPrefix))
	}
	if val, ok := r.Values[key]; ok {
		return val, true
	}
	return "", false
}

func (r *Request) getHeader(ctx context.Context, key string) (string, bool) {
	val, ok := r.Headers[key]
	if !ok {
		return "", false
	}
	if len(val) == 0 {
		return "", false
	}
	return val[0], true
}

func fillRequest(req *Request) error {
	md, ok := metadata.FromIncomingContext(req.Context)
	if !ok {
		return fmt.Errorf("fail to get metadata from incomming ctx")
	}
	logger.Infof(req.Context, "metadata: %+v", md)

	values := make(map[string]string)
	values[config.RequestFieldServiceName] = req.ServiceName
	values[config.RequestFieldMethodName] = req.MethodName

	fields := req.Msg.ProtoReflect().Descriptor().Fields()
	for idx := 0; idx < fields.Len(); idx++ {
		values[getFieldName(fields.Get(idx))] = req.Msg.ProtoReflect().Get(fields.Get(idx)).String()
	}
	req.Values = values
	req.Headers = md

	return nil
}

func getFieldName(desc protoreflect.FieldDescriptor) string {
	return fmt.Sprintf("body.%s", desc.Name())
}
