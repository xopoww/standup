package tg

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/logging"
)

type MockBot struct {
	mock.Mock
	uc chan tgbotapi.Update
}

func NewMockBot() *MockBot {
	return &MockBot{
		uc: make(chan tgbotapi.Update),
	}
}

func (b *MockBot) Updates() tgbotapi.UpdatesChannel {
	return b.uc
}

func (b *MockBot) Send(m tgbotapi.Chattable) (tgbotapi.Message, error) {
	result := b.Called(m)
	return result.Get(0).(tgbotapi.Message), result.Error(1)
}

func (b *MockBot) AddUpdate(ctx context.Context, t *testing.T, u tgbotapi.Update) {
	logging.L(ctx).Sugar().Debugf("Sending update to mock: %s...", logging.MarshalJSON(u))
	select {
	case b.uc <- u:
		return
	case <-ctx.Done():
		require.FailNow(t, "AddUpdate: %s.", ctx.Err())
	}
}

func (b *MockBot) Stop() {
	close(b.uc)
}
