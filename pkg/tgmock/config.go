package tgmock

type Config struct {
	Control string `yaml:"control" validate:"required,hostname_port"`
	Service string `yaml:"service" validate:"required,hostname_port"`
}
