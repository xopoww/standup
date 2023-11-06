package bot

import (
	"github.com/xopoww/standup/internal/bot/service"
	"github.com/xopoww/standup/internal/bot/tg"
)

const DefaultConfigPath = "/etc/standup/bot/config.yaml"

type Config struct {
	Bot *tg.Config `yaml:"bot" validate:"required"`

	Database *struct {
		DBS string `yaml:"dbs" validate:"required"`
	} `yaml:"database" validate:"required"`

	Service *service.Config `yaml:"service" validate:"required"`

	Standup *struct {
		Addr string `yaml:"addr" validate:"required,hostname_port"`
	} `yaml:"standup" validate:"required"`

	PrivateKeyFile string `yaml:"private_key_file" validate:"required,file"`
}
