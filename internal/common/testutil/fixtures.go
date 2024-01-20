package testutil

import (
	"context"
	"strconv"

	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/pkg/api/standup"
)

func OwnerID(ctx context.Context) int64 {
	id, err := strconv.ParseInt(TestID(ctx), 16, 64)
	if err != nil {
		logging.L(ctx).Sugar().Panicf("Unexpected TestID (%q): %s.", TestID(ctx), err)
	}
	return id
}

func CreateMessageRequest(ctx context.Context) *standup.CreateMessageRequest {
	return &standup.CreateMessageRequest{
		Text:    "Test message",
		OwnerId: OwnerID(ctx),
	}
}
