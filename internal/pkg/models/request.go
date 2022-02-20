package models

import (
	"context"
	"encoding/json"
)

const (
	FieldHost        = "Host"
	FieldServiceName = "ServiceName"
	FieldMethodName  = "MethodName"
)

// Request запрос в имитирующий сервис
type Request interface {
	GetValue(ctx context.Context, key string) (string, bool)
}

// Response ответ имитирующего сервиса
type Response struct {
	ID      string
	Message json.RawMessage
}
