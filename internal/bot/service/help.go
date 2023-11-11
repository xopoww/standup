package service

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/bot/formatting"
	"github.com/xopoww/standup/internal/bot/tg"
)

func (s *Service) help(_ context.Context, tm tgbotapi.Message) (err error) {
	if args := tm.CommandArguments(); args != "" {
		for _, cmd := range s.cmds {
			if cmd.Name != args {
				continue
			}
			reply := tg.NewMessagef(tm.Chat.ID, formatting.FormatCommandHelp(cmd))
			_, err = s.deps.Bot.Send(reply)
			if err != nil {
				err = fmt.Errorf("send reply: %w", err)
			}
			return err
		}
	}
	reply := tg.NewMessagef(tm.Chat.ID, formatting.FormatHelp(s.cmds))
	_, err = s.deps.Bot.Send(reply)
	if err != nil {
		err = fmt.Errorf("send reply: %w", err)
	}
	return err
}
