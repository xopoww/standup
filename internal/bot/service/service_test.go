package service_test

import (
	"context"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/service"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/logging"
)

func TestEcho(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	ctx = logging.WithLogger(ctx, logging.NewLogger())

	mb := tg.NewMockBot()
	defer mb.AssertExpectations(t)

	s, err := service.NewService(logging.L(ctx), mb)
	go s.Start()
	defer s.Stop()
	require.NoError(t, err)

	incoming := tgbotapi.Message{
		MessageID: 10,
		Chat: &tgbotapi.Chat{
			ID: 100,
		},
		Text: "Test text",
	}

	mb.On("Send", mock.MatchedBy(func(c tgbotapi.Chattable) bool {
		msg, ok := c.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}
		return msg.Text == incoming.Text &&
			msg.ChatID == incoming.Chat.ID &&
			msg.ReplyToMessageID == incoming.MessageID
	})).Return(tgbotapi.Message{MessageID: 20}, nil)

	mb.AddUpdate(ctx, t, tgbotapi.Update{Message: &incoming})
}
