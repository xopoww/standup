package testutil

import (
	"context"

	"github.com/xopoww/standup/pkg/api/standup"
)

func OwnerID(ctx context.Context) string {
	return "test-owner-" + TestID(ctx)
}

func CreateMessageRequest(ctx context.Context) *standup.CreateMessageRequest {
	return &standup.CreateMessageRequest{
		Text:    "Test message",
		OwnerId: OwnerID(ctx),
	}
}
