package tg

import "time"

type Config struct {
	TokenFile string `yaml:"token_file" validate:"required,file"`

	// TODO: add webhook option
	Poll *struct {
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"poll" validate:"required"`
}

const (
	defaultPollTimeout = time.Minute
)

func (c *Config) SetDefaults() {
	if c.Poll != nil && c.Poll.Timeout == 0 {
		c.Poll.Timeout = defaultPollTimeout
	}
}
