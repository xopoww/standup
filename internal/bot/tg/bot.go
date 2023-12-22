package tg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/internal/common/logging"
)

type bot struct {
	cfg Config

	tb *tgbotapi.BotAPI
	uc tgbotapi.UpdatesChannel
}

func NewBot(ctx context.Context, cfg Config, devel bool) (Bot, error) {
	cfg.SetDefaults()

	data, err := os.ReadFile(cfg.TokenFile)
	if err != nil {
		return nil, fmt.Errorf("read token: %w", err)
	}
	token := strings.TrimSpace(string(data))

	client := &http.Client{}
	if cfg.HTTPLogging {
		client.Transport = logging.RoundTripper(logging.L(ctx), http.DefaultTransport)
	}
	tb, err := tgbotapi.NewBotAPIWithClient(token, cfg.APIEndpoint, client)
	if err != nil {
		return nil, fmt.Errorf("new bot api: %w", err)
	}
	tb.Debug = devel
	logging.L(ctx).Sugar().Infof("Authorized as %q.", tb.Self.UserName)

	b := &bot{cfg: cfg, tb: tb}

	switch {
	case cfg.Poll != nil:
		b.initLongPolling(ctx)
		err = nil
	default:
		err = errors.New("invalid config: no update method")
	}
	if err != nil {
		return nil, fmt.Errorf("init updates: %w", err)
	}

	return b, nil
}

func (b *bot) initLongPolling(ctx context.Context) {
	uCfg := tgbotapi.NewUpdate(0)
	uCfg.Timeout = int(b.cfg.Poll.Timeout / time.Second)
	logging.L(ctx).Sugar().Infof("Starting long polling (t/o %s) for updates...", b.cfg.Poll.Timeout)
	b.uc = b.tb.GetUpdatesChan(uCfg)
}

func (b *bot) Updates() tgbotapi.UpdatesChannel {
	return b.uc
}

func (b *bot) Send(m tgbotapi.Chattable) (tgbotapi.Message, error) {
	return b.tb.Send(m)
}

func (b *bot) Stop() {
	b.tb.StopReceivingUpdates()
}
