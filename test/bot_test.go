package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/pkg/tgmock/tests"
)

func TestBotWhitelist(t *testing.T) {
	cc := []struct {
		name        string
		whitelisted bool
	}{
		{
			name:        "allowed",
			whitelisted: true,
		},
		{
			name:        "restricted",
			whitelisted: false,
		},
	}
	for _, c := range cc {
		RunTest(t, c.name, func(ctx context.Context, t *testing.T) {
			user := tests.ContextUser(ctx)
			chat := tests.ContextChat(ctx)

			if c.whitelisted {
				require.NoError(t, deps.Repo.SetWhitelisted(ctx, user.GetId(), true))
			}

			mid := tests.SendMessage(ctx, t, deps.TM, user, chat, "/start")

			msg := tests.WaitForMessagesTimeout(ctx, t, deps.TM, chat, 1, time.Minute)[0]
			if c.whitelisted {
				require.NotContains(t, msg.GetText(), "currently restricted")
			} else {
				require.Contains(t, msg.GetText(), "currently restricted")
				require.Equal(t, mid, msg.GetReplyToMessage().GetMessageId())
			}

			logging.L(ctx).Sugar().Debugf("Chat %d history:\n%s", chat.GetId(), tests.ChatHistory(ctx, t, deps.TM, chat))
		})
	}
}
