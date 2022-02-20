package logger

import (
	"context"
	"fmt"
	"log"
)

func Info(ctx context.Context, msg string) {
	log.Println(msg)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	log.Println(fmt.Sprintf(msg, args))
}

func Fatal(ctx context.Context, msg string) {
	log.Fatalln(msg)
}

func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	log.Fatalln(fmt.Sprintf(msg, args))
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	log.Println(fmt.Sprintf(msg, args))
}

func Error(ctx context.Context, msg string) {
	log.Println(msg)
}
