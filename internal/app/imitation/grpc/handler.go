package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Speakerkfm/iso/internal/pkg/models"
	"github.com/Speakerkfm/iso/internal/pkg/request_processor"
)

type handler struct {
	processor request_processor.Processor

	svc map[string]Service
}

func NewHandler(processor request_processor.Processor) *handler {
	return &handler{
		processor: processor,
		svc:       make(map[string]Service),
	}
}

func (h *handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	svc, isRegistered := h.svc[req.ServiceName]
	if !isRegistered {
		return nil, fmt.Errorf("service: %s not registered", req.ServiceName)
	}

	method, isRegistered := svc.Methods[req.MethodName]
	if !isRegistered {
		return nil, fmt.Errorf("method: %s not registered", req.MethodName)
	}

	// convert grpc request to model

	resp, err := h.processor.Process(ctx, &models.Request{})
	if err != nil {
		return nil, fmt.Errorf("fail to process request: %s", err.Error())
	}

	msg := method.RespStruct.ProtoReflect().New().Interface()
	if err := json.Unmarshal(resp.Message, &msg); err != nil {
		return nil, fmt.Errorf("fail to unmarshal resp json into proto struct")
	}

	protoResp := &Response{
		msg: msg,
	}

	// save in store by message id

	return protoResp, nil
}