package tg

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewReplyf(to tgbotapi.Message, template string, args ...any) tgbotapi.MessageConfig {
	cfg := tgbotapi.NewMessage(to.Chat.ID, fmt.Sprintf(template, args...))
	cfg.ReplyToMessageID = to.MessageID
	return cfg
}
