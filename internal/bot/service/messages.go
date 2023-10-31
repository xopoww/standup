package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/bot/formatting"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) getReport(ctx context.Context, tm tgbotapi.Message) (err error) {
	now := tm.Time().UTC()

	args := strings.Split(tm.CommandArguments(), " ")
	if len(args) < 1 || len(args) > 3 {
		return NewSyntaxError("bad number of arguments")
	}

	from, err := formatting.ParseTime(args[0], now)
	if err != nil {
		return NewSyntaxError(err.Error())
	}
	var to time.Time
	if len(args) == 2 {
		to, err = formatting.ParseTime(args[1], now)
		if err != nil {
			return NewSyntaxError(err.Error())
		}
	} else {
		to = now
	}

	ctx, err = s.issueToken(ctx, tm.From.UserName, ShortTTl)
	if err != nil {
		return fmt.Errorf("issue token: %w", err)
	}

	rsp, err := s.deps.Client.ListMessages(ctx, &standup.ListMessagesRequest{
		OwnerId: tm.From.UserName,
		From:    timestamppb.New(from),
		To:      timestamppb.New(to),
	})
	if err != nil {
		return fmt.Errorf("list messages: %w", err)
	}

	reply := tgbotapi.NewMessage(tm.Chat.ID, formatting.FormatMessages("Report", rsp.GetMessages()))
	reply.ParseMode = "MarkdownV2"
	_, err = s.deps.Bot.Send(reply)
	if err != nil {
		return fmt.Errorf("send reply: %w", err)
	}
	return nil
}

func (s *Service) addMessage(ctx context.Context, tm tgbotapi.Message) (err error) {
	if tm.Text == "" {
		return nil
	}

	ctx, err = s.issueToken(ctx, tm.From.UserName, ShortTTl)
	if err != nil {
		return fmt.Errorf("issue token: %w", err)
	}

	rsp, err := s.deps.Client.CreateMessage(ctx, &standup.CreateMessageRequest{
		Text:    tm.Text,
		OwnerId: tm.From.UserName,
	})
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	_, err = s.deps.Bot.Send(tg.NewReplyf(tm, "Created message %q.", rsp.GetId()))
	if err != nil {
		return fmt.Errorf("send reply: %w", err)
	}
	return nil
}
