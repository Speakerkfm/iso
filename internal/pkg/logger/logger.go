package logger

import (
	"context"
	"log"

	"github.com/Speakerkfm/iso/internal/pkg/config"
)

func Info(ctx context.Context, msg string) {
	if config.LoggerLevel > 1 {
		return
	}
	log.Println(msg)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	if config.LoggerLevel > 1 {
		return
	}
	log.Printf(msg, args...)
}

func Warn(ctx context.Context, msg string) {
	if config.LoggerLevel > 2 {
		return
	}
	log.Println(msg)
}

func Warnf(ctx context.Context, msg string, args ...interface{}) {
	if config.LoggerLevel > 2 {
		return
	}
	log.Printf(msg, args...)
}

func Fatal(ctx context.Context, msg string) {
	log.Fatalln(msg)
}

func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

func Error(ctx context.Context, msg string) {
	log.Println(msg)
}
