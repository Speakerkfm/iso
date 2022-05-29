package models

import (
	"context"
	"encoding/json"
	"time"
)

// Request запрос в имитирующий сервис
type Request interface {
	GetValue(ctx context.Context, key string) (string, bool)
	GetHandledAt() time.Time
}

// Response ответ имитирующего сервиса
type Response struct {
	ID      string
	Message json.RawMessage
	Error   string
}
