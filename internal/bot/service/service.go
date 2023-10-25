package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/logging"
	"go.uber.org/zap"
)

type Service struct {
	b      tg.Bot
	logger *zap.Logger
	wg     sync.WaitGroup
}

func NewService(logger *zap.Logger, bot tg.Bot) (*Service, error) {
	return &Service{b: bot, logger: logger}, nil
}

func (s *Service) Start() {
	for i := 0; i < 1; i++ {
		s.wg.Add(1)
		go s.worker(s.b.Updates())
	}
}

func (s *Service) Stop() {
	s.b.Stop()
	s.wg.Wait()
}

func (s *Service) worker(jobs <-chan tgbotapi.Update) {
	defer s.wg.Done()
	for u := range jobs {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			ctx = logging.WithLogger(ctx, s.logger.With(
				zap.String("upd", fmt.Sprintf("%08d", u.UpdateID)),
			))

			logging.L(ctx).Sugar().Debugf("Update received: %s.", logging.MarshalJSON(u))

			start := time.Now()
			err := s.handleUpdate(ctx, u)
			delta := time.Since(start)

			if err != nil {
				logging.L(ctx).Sugar().Errorf("Failed to handle update (%s): %s.", delta, err)
			} else {
				logging.L(ctx).Sugar().Debugf("Update handled (%s).", delta)
			}
		}()
	}
}

func (s *Service) handleUpdate(ctx context.Context, u tgbotapi.Update) error {
	if u.Message == nil {
		return nil
	}
	msg := u.Message
	logging.L(ctx).Sugar().Debugf("Got message %q.", msg.Text)

	reply := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
	reply.ReplyToMessageID = msg.MessageID

	sent, err := s.b.Send(reply)
	if err != nil {
		return fmt.Errorf("send reply: %w", err)
	}
	logging.L(ctx).Sugar().Debugf("Sent reply with id %d.", sent.MessageID)
	return nil
}
