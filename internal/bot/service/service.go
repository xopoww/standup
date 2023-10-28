package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/auth"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/logging"
	"github.com/xopoww/standup/pkg/api/standup"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	b      tg.Bot
	logger *zap.Logger
	wg     sync.WaitGroup

	sc standup.StandupClient
	i  auth.Issuer
}

func NewService(logger *zap.Logger, bot tg.Bot, sc standup.StandupClient, i auth.Issuer) (*Service, error) {
	return &Service{b: bot, logger: logger, sc: sc, i: i}, nil
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
		logging.L(ctx).Sugar().Debug("Skip unsupported update.")
		return nil
	}
	msg := *u.Message

	var err error
	switch cmd := msg.Command(); cmd {
	case "":
		err = s.addMessage(ctx, msg)
	case "report":
		err = s.getReport(ctx, msg)
	default:
		err = NewSyntaxErrorf("unknown command: %q", cmd)
	}
	if err == nil {
		return nil
	}

	serr := SyntaxError{}
	if errors.As(err, &serr) {
		_, rerr := s.b.Send(tg.NewReplyf(msg, "Invalid command: %s.", serr.Error()))
		return rerr
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.PermissionDenied:
			_, rerr := s.b.Send(tg.NewReplyf(msg, "Permission denied."))
			return rerr
		}
	}

	_, rerr := s.b.Send(tg.NewReplyf(msg, "Internal error occured."))
	if rerr != nil {
		err = fmt.Errorf("send reply: %w (while handling error: %s)", rerr, err)
	}
	return err
}
