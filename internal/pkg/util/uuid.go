package util

import (
	uuid "github.com/satori/go.uuid"
)

// NewUUID генерирует новый uuid
func NewUUID() string {
	return uuid.NewV4().String()
}
