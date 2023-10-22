package logging

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	l, err := cfg.Build(zap.AddCaller())
	if err != nil {
		log.Fatalf("Failed to init logging: %s.", err)
	}
	return l
}
