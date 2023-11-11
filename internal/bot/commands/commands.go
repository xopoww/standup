package commands

import (
	_ "embed"
	"io"
	"strings"

	"github.com/xopoww/standup/pkg/config"
)

//go:embed commands.yaml
var commandsData string

type Desc struct {
	Name  string `yaml:"name" validate:"required"`
	Usage string `yaml:"usage"`
	Short string `yaml:"short" validate:"required"`
	Long  string `yaml:"long"`
}

func (d Desc) sanitized() Desc {
	d.Usage = strings.TrimSpace(d.Usage)
	d.Short = strings.TrimSpace(d.Short)
	d.Long = strings.TrimSpace(d.Long)
	return d
}

func LoadDescriptions() ([]Desc, error) {
	return loadFrom(strings.NewReader(commandsData))
}

func loadFrom(r io.Reader) ([]Desc, error) {
	var cmds struct {
		Commands []Desc `yaml:"commands" validate:"required"`
	}
	err := config.Load(r, &cmds)
	for i := range cmds.Commands {
		cmds.Commands[i] = cmds.Commands[i].sanitized()
	}
	return cmds.Commands, err
}
