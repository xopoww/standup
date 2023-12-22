package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/pkg/tgmock/control"
)

func TestBot(t *testing.T) {
	RunTest(t, "default", func(ctx context.Context, t *testing.T) {
		_, err := deps.TM.CreateUpdate(ctx, &control.CreateUpdateRequest{
			Update: &control.Update{
				Message: &control.Message{
					From: &control.User{
						Id:       12345,
						Username: "test-user",
					},
					Date: time.Now().Unix(),
					Chat: &control.Chat{
						Id: 12345,
					},
					Text: "/start",
				},
			},
		})
		require.NoError(t, err)

		time.Sleep(time.Second)

		rsp, err := deps.TM.ListMessages(ctx, &control.ListMessagesRequest{
			ChatId: 12345,
		})
		require.NoError(t, err)
		require.Len(t, rsp.GetMessages(), 1)
	})
}
