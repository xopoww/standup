package daemon

type Config struct {
	Database struct {
		DBS string `yaml:"dbs" validate:"required"`
	} `yaml:"database" validate:"required"`

	Service struct {
		Addr string `yaml:"addr" validate:"required,hostname_port"`
	} `yaml:"service" validate:"required"`
}
