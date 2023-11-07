//nolint:lll // RunTest signature
package service_test

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/common/auth"
	"github.com/xopoww/standup/internal/common/testutil"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc/metadata"
)

func TestAddMessage(t *testing.T) {
	RunTest("default", t, func(ctx context.Context, t *testing.T, bot tg.MockBot, bc tg.MockBotClient, sc *testutil.MockStandupClient) {
		const text = "Test text of the message."
		const id = "01234567890abcde"

		sc.On("CreateMessage", testutil.OutgoingMetadata(func(md metadata.MD) bool {
			v := md.Get(auth.GRPCMetadataTokenKey)
			return len(v) > 0 && v[0] == TestUserName+"_token"
		}), mock.MatchedBy(func(req *standup.CreateMessageRequest) bool {
			return req.GetOwnerId() == TestUserName &&
				req.GetText() == text
		})).Return(&standup.CreateMessageResponse{Id: id}, nil)

		msg := NewIncomingMessage(text)
		bc.SendMessage(ctx, t, msg)
		reply := bc.RecvMessage(ctx, t)
		require.Equal(t, msg.Chat.ID, reply.ChatID)
		require.Equal(t, msg.MessageID, reply.ReplyToMessageID)
		require.Contains(t, reply.Text, id)
	})

	RunTest("not_textual", t, func(ctx context.Context, t *testing.T, bot tg.MockBot, bc tg.MockBotClient, sc *testutil.MockStandupClient) {
		msg := NewIncomingMessage("")
		msg.Sticker = &tgbotapi.Sticker{}

		bc.SendMessage(ctx, t, msg)
	})
}
