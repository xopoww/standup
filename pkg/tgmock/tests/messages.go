package tests

import (
	"context"
	"errors"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/pkg/tgmock/control"
)

func SendMessage(ctx context.Context, t require.TestingT,
	tm control.TGMockControlClient, u *control.User, c *control.Chat, text string) int64 {
	rsp, err := tm.CreateUpdate(ctx, &control.CreateUpdateRequest{
		Update: &control.Update{
			Message: &control.Message{
				From: u,
				Date: time.Now().Unix(),
				Chat: c,
				Text: text,
			},
		},
	})
	require.NoError(t, err)
	return rsp.GetMessageId()
}

func WaitForMessagesTimeout(ctx context.Context, t require.TestingT,
	tm control.TGMockControlClient, c *control.Chat, n int, timeout time.Duration) []*control.Message {
	ctx, cancel := context.WithTimeoutCause(ctx, timeout, errors.New("wait for messages timed out"))
	defer cancel()

	const pollPeriod = time.Millisecond * 100
	tck := time.NewTicker(pollPeriod)
	for {
		select {
		case <-ctx.Done():
			require.Fail(t, ctx.Err().Error())
			return nil
		case <-tck.C:
			rsp, err := tm.ListMessages(ctx, &control.ListMessagesRequest{ChatId: c.GetId()})
			require.NoError(t, err)
			if len(rsp.GetMessages()) >= n {
				return rsp.GetMessages()
			}
		}
	}
}
