package models

import (
	"encoding/json"
)

// Request запрос в имитирующий сервис
type Request struct {
}

// Response ответ имитирующего сервиса
type Response struct {
	ID      string
	Message json.RawMessage
}
