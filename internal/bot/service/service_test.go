package service_test

import (
	"context"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/service"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/logging"
)

func TestEcho(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	ctx = logging.WithLogger(ctx, logging.NewLogger())

	bot, client := tg.NewMockBot()
	s, err := service.NewService(logging.L(ctx), bot)
	require.NoError(t, err)

	go s.Start()
	defer s.Stop()

	incoming := tgbotapi.Message{
		MessageID: 10,
		Chat: &tgbotapi.Chat{
			ID: 100,
		},
		Text: "Test text",
	}

	client.SendMessage(ctx, t, incoming)
	reply := client.RecvMessage(ctx, t)
	require.Equal(t, incoming.Text, reply.Text)
	require.Equal(t, incoming.Chat.ID, reply.ChatID)
	require.Equal(t, incoming.MessageID, reply.ReplyToMessageID)
}
