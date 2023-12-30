package tg

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	APIEndpoint string `yaml:"api_endpoint"`

	TokenFile string `yaml:"token_file" validate:"required,file"`

	// TODO: add webhook option
	Poll *struct {
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"poll" validate:"required"`

	HTTPLogging bool `yaml:"http_logging"`
}

const (
	defaultPollTimeout = time.Minute
)

func (c *Config) SetDefaults() {
	if c.Poll != nil && c.Poll.Timeout == 0 {
		c.Poll.Timeout = defaultPollTimeout
	}
	if c.APIEndpoint == "" {
		c.APIEndpoint = tgbotapi.APIEndpoint
	}
}
