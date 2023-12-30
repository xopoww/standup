package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/bot/commands"
	"github.com/xopoww/standup/internal/bot/commands/commandtypes"
	"github.com/xopoww/standup/internal/bot/formatting"
	"github.com/xopoww/standup/internal/bot/models"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/common/auth"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/internal/common/repository/dberrors"
	"github.com/xopoww/standup/pkg/api/standup"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Deps struct {
	Bot    tg.Bot
	Models models.Models
	Client standup.StandupClient
	Issuer auth.Issuer
}

type Service struct {
	logger *zap.Logger
	wg     sync.WaitGroup

	cfg  Config
	deps Deps

	cmds []commandtypes.Desc
}

func NewService(logger *zap.Logger, cfg Config, deps Deps) (*Service, error) {
	cmds, err := commands.LoadDescriptions()
	if err != nil {
		return nil, fmt.Errorf("load descriptions: %w", err)
	}
	logger.Sugar().Debugf("Loaded descriptions of %d command(s).", len(cmds))
	return &Service{
		logger: logger,
		cfg:    cfg,
		deps:   deps,
		cmds:   cmds,
	}, nil
}

func (s *Service) Start() {
	for i := 0; i < 1; i++ {
		s.wg.Add(1)
		go s.worker(s.deps.Bot.Updates())
	}
}

func (s *Service) Stop() {
	s.deps.Bot.Stop()
	s.wg.Wait()
}

func (s *Service) worker(jobs <-chan tgbotapi.Update) {
	const updateTimeout = 10 * time.Second
	defer s.wg.Done()
	for u := range jobs {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
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

	if allowed, err := s.handleMessageSender(ctx, msg); err != nil {
		// do not report errs to unverified users
		return err
	} else if !allowed {
		return nil
	}

	var err error
	switch cmd := msg.Command(); cmd {
	case "":
		err = s.addMessage(ctx, msg)
	case "report":
		err = s.getReport(ctx, msg)
	case "help":
		err = s.help(ctx, msg)
	case "start":
		err = s.start(ctx, msg)
	default:
		err = formatting.NewSyntaxErrorf("unknown command %q", cmd)
	}
	if err == nil {
		return nil
	}

	// TODO: add command to syntax errors
	//nolint:errorlint // Printable is not an error
	if perr, ok := err.(formatting.Printable); ok {
		_, rerr := s.deps.Bot.Send(tg.NewReplyf(msg, perr.Printable()))
		return rerr
	}

	if st, ok := status.FromError(err); ok {
		if st.Code() == codes.PermissionDenied {
			_, rerr := s.deps.Bot.Send(tg.NewReplyf(msg, formatting.PermissionDenied))
			return rerr
		}
	}

	_, rerr := s.deps.Bot.Send(tg.NewReplyf(msg, formatting.InternalError))
	if rerr != nil {
		err = fmt.Errorf("send reply: %w (while handling error: %w)", rerr, err)
	}
	return err
}

func (s *Service) handleMessageSender(ctx context.Context, msg tgbotapi.Message) (bool, error) {
	user := models.FromTG(msg.From)

	// TODO: rm after transition period
	if _, err := s.deps.Models.GetUserByID(ctx, user.ID); errors.Is(err, dberrors.ErrNotFound) {
		err := s.deps.Models.SetUserID(ctx, user.Username, user.ID)
		if err == nil {
			logging.L(ctx).Sugar().Infof("Set id for user %s.", user)
		} else if !errors.Is(err, dberrors.ErrNotFound) {
			return false, fmt.Errorf("set user id: %w", err)
		}
	}

	if err := s.deps.Models.UpsertUser(ctx, user); err != nil {
		return false, fmt.Errorf("upsert user: %w", err)
	}

	if !s.cfg.WhitelistEnabled {
		return true, nil
	}
	allowed, err := s.deps.Models.GetWhitelisted(ctx, user.ID)
	if err != nil {
		return false, fmt.Errorf("check whitelist: %w", err)
	}

	if allowed {
		return true, nil
	}
	logging.L(ctx).Sugar().Debugf("User %s is not in whitelist, access denied.", user)
	_, err = s.deps.Bot.Send(tg.NewReplyf(msg, formatting.UsageRestricted))
	return false, err
}
