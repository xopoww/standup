package service_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/service"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/internal/common/testutil"
)

type testFunc func(
	ctx context.Context, t *testing.T, bot tg.MockBot,
	bc tg.MockBotClient, sc *testutil.MockStandupClient,
)

func RunTest(name string, t *testing.T, f testFunc) {
	t.Run(name, func(tt *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		logger := logging.NewLogger()
		defer func() {
			_ = logger.Sync()
		}()
		ctx = logging.WithLogger(ctx, logger)

		sc := &testutil.MockStandupClient{}
		defer sc.AssertExpectations(tt)

		bot, bc := tg.NewMockBot()
		defer bc.RequireEmpty(tt)
		cfg := service.Config{WhitelistEnabled: false}
		srv, err := service.NewService(logging.L(ctx), cfg, service.Deps{
			Bot:    bot,
			Models: nil,
			Client: sc,
			Issuer: TestIssuer{},
		})
		require.NoError(t, err)

		srv.Start()
		defer srv.Stop()
		f(ctx, tt, bot, bc, sc)
	})
}

const TestUserName = "test_user"

func NewIncomingMessage(text string) tgbotapi.Message {
	return tgbotapi.Message{
		MessageID: rand.Int(),
		From: &tgbotapi.User{
			ID:        7357,
			IsBot:     false,
			FirstName: "Test",
			LastName:  "User",
			UserName:  TestUserName,
		},
		Date: int(time.Now().Unix()),
		Chat: &tgbotapi.Chat{
			ID:   73570,
			Type: "private",
		},
		Text: text,
	}
}

type TestIssuer struct{}

func (i TestIssuer) IssueToken(subjectID string, _, _ time.Time) (string, error) {
	return subjectID + "_token", nil
}
