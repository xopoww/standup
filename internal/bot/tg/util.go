package tg

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/bot/formatting"
)

func NewMessagef(chatID int64, format string, a ...any) tgbotapi.MessageConfig {
	cfg := tgbotapi.NewMessage(chatID, fmt.Sprintf(format, a...))
	cfg.ParseMode = formatting.ParseMode
	return cfg
}

func NewReplyf(to tgbotapi.Message, format string, a ...any) tgbotapi.MessageConfig {
	cfg := NewMessagef(to.Chat.ID, format, a...)
	cfg.ReplyToMessageID = to.MessageID
	return cfg
}
