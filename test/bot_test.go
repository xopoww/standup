package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/internal/common/repository/dberrors"
	"github.com/xopoww/standup/pkg/tgmock/tests"
)

func TestBot(t *testing.T) {
	RunTest(t, "whitelist", func(ctx context.Context, t *testing.T) {
		u := tests.ContextUser(ctx)
		err := deps.Repo.SetWhitelisted(ctx, u.GetId(), true)
		require.ErrorIs(t, err, dberrors.ErrNotFound)

		c := tests.ContextChat(ctx)

		mid := tests.SendMessage(ctx, t, deps.TM, u, c, "/start")

		msgs := tests.WaitForMessagesTimeout(ctx, t, deps.TM, c, 1, time.Minute)
		require.Contains(t, msgs[0].GetText(), "currently restricted")
		require.Equal(t, mid, msgs[0].GetReplyToMessage().GetMessageId())

		err = deps.Repo.SetWhitelisted(ctx, u.GetId(), true)
		require.NoError(t, err)

		tests.SendMessage(ctx, t, deps.TM, u, c, "/start")

		msgs = tests.WaitForMessagesTimeout(ctx, t, deps.TM, c, 2, time.Minute)
		require.NotContains(t, msgs[1].GetText(), "currently restricted")

		logging.L(ctx).Sugar().Debugf("Chat %d history:\n%s", c.GetId(), tests.ChatHistory(ctx, t, deps.TM, c))
	})
}
