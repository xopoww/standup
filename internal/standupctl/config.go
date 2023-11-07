package standupctl

const DefaultConfigPath = "/etc/standup/standupctl/config.yaml"

type Config struct {
	Database struct {
		DBS string `yaml:"dbs" validate:"required"`
	} `yaml:"database" validate:"required"`
}
