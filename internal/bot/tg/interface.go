package tg

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Bot interface {
	Updates() tgbotapi.UpdatesChannel
	Send(m tgbotapi.Chattable) (tgbotapi.Message, error)
	Stop()
}
