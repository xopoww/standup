package testutil

import (
	"context"
	"fmt"

	"github.com/xopoww/standup/pkg/identifiers"
)

type testIDKey struct{}

func NewContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	id, err := identifiers.GenerateShortID()
	if err != nil {
		panic(fmt.Sprintf("standup: cound not generate test id: %s", err))
	}
	ctx = context.WithValue(ctx, testIDKey{}, id)

	return ctx, cancel
}

func TestID(ctx context.Context) string {
	id, ok := ctx.Value(testIDKey{}).(string)
	if !ok {
		return "unknown"
	}
	return id
}
