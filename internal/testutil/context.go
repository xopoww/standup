package testutil

import (
	"context"

	"github.com/xopoww/standup/pkg/identifiers"
)

var testIdKey struct{}

func NewContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	id := identifiers.GenerateShortID()
	ctx = context.WithValue(ctx, testIdKey, id)

	return ctx, cancel
}

func TestID(ctx context.Context) string {
	return ctx.Value(testIdKey).(string)
}
