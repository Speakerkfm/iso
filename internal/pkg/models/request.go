package models

import (
	"encoding/json"
)

const (
	FieldHost        = "Host"
	FieldServiceName = "ServiceName"
	FieldMethodName  = "MethodName"
)

// Request запрос в имитирующий сервис
type Request struct {
	Values map[string]string
}

// Response ответ имитирующего сервиса
type Response struct {
	ID      string
	Message json.RawMessage
}
