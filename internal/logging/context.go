package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func L(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return l
}
