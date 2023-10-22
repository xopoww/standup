package testutil

import (
	"context"

	"github.com/xopoww/standup/pkg/identifiers"
)

type testIdKey struct{}

func NewContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	id := identifiers.GenerateShortID()
	ctx = context.WithValue(ctx, testIdKey{}, id)

	return ctx, cancel
}

func TestID(ctx context.Context) string {
	id, ok := ctx.Value(testIdKey{}).(string)
	if !ok {
		return "unknown"
	}
	return id
}
