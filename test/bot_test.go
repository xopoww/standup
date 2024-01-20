package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/formatting"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/tgmock/tests"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		RunBotTest(t, c.name, func(ctx context.Context, t *testing.T) {
			chat := tests.ContextChat(ctx)

			user := tests.ContextUser(ctx)
			if c.whitelisted {
				require.NoError(t, deps.Repo.SetWhitelisted(ctx, user.GetId(), true))
			}

			mid := tests.SendMessage(ctx, t, deps.TM, user, chat, "/start")

			msg := tests.WaitForMessages(ctx, t, deps.TM, chat, 1, time.Minute)[0]
			if c.whitelisted {
				require.NotContains(t, msg.GetText(), "currently restricted")
			} else {
				require.Contains(t, msg.GetText(), "currently restricted")
				require.Equal(t, mid, msg.GetReplyToMessage().GetMessageId())
			}
		})
	}
}

func TestBotHelp(t *testing.T) {
	RunBotTest(t, "default", func(ctx context.Context, t *testing.T) {
		chat := tests.ContextChat(ctx)

		user := tests.ContextUser(ctx)
		require.NoError(t, deps.Repo.SetWhitelisted(ctx, user.GetId(), true))

		// basic help
		tests.SendMessage(ctx, t, deps.TM, user, chat, "/help")
		msg := tests.WaitForMessages(ctx, t, deps.TM, chat, 1, time.Minute)[0]
		require.Contains(t, msg.GetText(), "Available commands")

		// help about one command
		tests.SendMessage(ctx, t, deps.TM, user, chat, "/help report")
		msg = tests.WaitForMessagesSince(ctx, t, deps.TM, chat, 1, time.Minute, msg.GetMessageId())[0]
		require.Contains(t, msg.GetText(), "Get messages")

		// help about unexisting command
		tests.SendMessage(ctx, t, deps.TM, user, chat, "/help foobar")
		msg = tests.WaitForMessagesSince(ctx, t, deps.TM, chat, 1, time.Minute, msg.GetMessageId())[0]
		require.Contains(t, msg.GetText(), "Available commands")
	})
}

func TestBotCreateMessage(t *testing.T) {
	RunBotTest(t, "default", func(ctx context.Context, t *testing.T) {
		chat := tests.ContextChat(ctx)

		user := tests.ContextUser(ctx)
		require.NoError(t, deps.Repo.SetWhitelisted(ctx, user.GetId(), true))

		const text = "Some text message."
		from := time.Now()
		tests.SendMessage(ctx, t, deps.TM, user, chat, text)
		msg := tests.WaitForMessages(ctx, t, deps.TM, chat, 1, time.Minute)[0]
		require.Contains(t, msg.GetText(), "Created message")

		rsp, err := deps.Client.ListMessages(withToken(ctx, t, user.GetId()),
			&standup.ListMessagesRequest{
				OwnerId: user.GetId(),
				From:    timestamppb.New(from),
				To:      timestamppb.Now(),
			},
		)
		require.NoError(t, err)
		require.Len(t, rsp.GetMessages(), 1)
		require.Equal(t, text, rsp.GetMessages()[0].GetText())
		require.Contains(t, msg.GetText(), rsp.GetMessages()[0].GetId())
	})
}

func TestBotReport(t *testing.T) {
	RunBotTest(t, "default", func(ctx context.Context, t *testing.T) {
		chat := tests.ContextChat(ctx)

		user := tests.ContextUser(ctx)
		require.NoError(t, deps.Repo.SetWhitelisted(ctx, user.GetId(), true))

		const text = "Some text message."
		tests.SendMessage(ctx, t, deps.TM, user, chat, text)
		msg := tests.WaitForMessages(ctx, t, deps.TM, chat, 1, time.Minute)[0]
		require.Contains(t, msg.GetText(), "Created message")

		tests.SendMessage(ctx, t, deps.TM, user, chat, "/report -1d")
		msg = tests.WaitForMessagesSince(ctx, t, deps.TM, chat, 1, time.Minute, msg.GetMessageId())[0]
		require.Contains(t, msg.GetText(), "Report")
		require.Contains(t, msg.GetText(), formatting.Escape(text))
	})
}
