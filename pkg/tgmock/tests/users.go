package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/xopoww/standup/internal/common/testutil"
	"github.com/xopoww/standup/pkg/tgmock/control"
)

func GenerateID() (int64, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, err
	}
	return id.Int64(), nil
}

func UserFromID(ctx context.Context, id int64) *control.User {
	tid := testutil.TestID(ctx)
	return &control.User{
		Id:       id,
		Username: fmt.Sprintf("u-%s-%x", tid, uint8(id)),
	}
}

func ChatFromID(_ context.Context, id int64) *control.Chat {
	return &control.Chat{
		Id: id,
	}
}

type defaultUserIDKey struct{}

func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, defaultUserIDKey{}, id)
}

func ContextUser(ctx context.Context) *control.User {
	id, ok := ctx.Value(defaultUserIDKey{}).(int64)
	if !ok {
		panic("tgmock: no default test user in context")
	}
	return UserFromID(ctx, id)
}

type defaultChatIDKey struct{}

func WithChatID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, defaultChatIDKey{}, id)
}

func ContextChat(ctx context.Context) *control.Chat {
	id, ok := ctx.Value(defaultChatIDKey{}).(int64)
	if !ok {
		panic("tgmock: no default test chat in context")
	}
	return ChatFromID(ctx, id)
}
