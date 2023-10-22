package logging

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	l, err := zap.NewDevelopmentConfig().Build(zap.AddCaller())
	if err != nil {
		log.Fatalf("Failed to init logging: %s.", err)
	}
	return l
}
