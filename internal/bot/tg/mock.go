package tg

import (
	"context"
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/common/logging"
)

type MockBot struct {
	cc chan tgbotapi.Chattable
	uc chan tgbotapi.Update
}

func NewMockBot() (MockBot, MockBotClient) {
	cc := make(chan tgbotapi.Chattable, 1)
	uc := make(chan tgbotapi.Update)
	return MockBot{cc, uc}, MockBotClient{cc, uc}
}

func (b MockBot) Updates() tgbotapi.UpdatesChannel {
	return b.uc
}

func (b MockBot) Send(m tgbotapi.Chattable) (tgbotapi.Message, error) {
	msg, ok := m.(tgbotapi.MessageConfig)
	if !ok {
		return tgbotapi.Message{}, fmt.Errorf("unsupported chattable: %+v", m)
	}
	b.cc <- m
	return tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: msg.ChatID,
		},
		Text: msg.Text,
		ReplyToMessage: &tgbotapi.Message{
			MessageID: msg.ReplyToMessageID,
		},
	}, nil
}

func (b MockBot) Stop() {
	close(b.uc)
}

type MockBotClient struct {
	cc <-chan tgbotapi.Chattable
	uc chan<- tgbotapi.Update
}

func (c MockBotClient) Recv(ctx context.Context, t *testing.T) tgbotapi.Chattable {
	select {
	case m := <-c.cc:
		return m
	case <-ctx.Done():
		require.FailNow(t, "MockBotClient.Recv: %s.", ctx.Err())
		return nil
	}
}

func (c MockBotClient) RecvMessage(ctx context.Context, t *testing.T) tgbotapi.MessageConfig {
	msg, ok := c.Recv(ctx, t).(tgbotapi.MessageConfig)
	require.True(t, ok, "Wrong tgbotapi.Chattable received")
	return msg
}

func (c MockBotClient) RequireEmpty(t *testing.T) {
	require.Zero(t, len(c.cc), "MockBot should have no outgoing messages")
}

func (c MockBotClient) SendUpdate(ctx context.Context, t *testing.T, u tgbotapi.Update) {
	logging.L(ctx).Sugar().Debugf("Sending update to mock: %s...", logging.MarshalJSON(u))
	select {
	case c.uc <- u:
		return
	case <-ctx.Done():
		require.FailNow(t, "MockBotClient.SendUpdate: %s.", ctx.Err())
	}
}

func (c MockBotClient) SendMessage(ctx context.Context, t *testing.T, msg tgbotapi.Message) {
	c.SendUpdate(ctx, t, tgbotapi.Update{Message: &msg})
}
