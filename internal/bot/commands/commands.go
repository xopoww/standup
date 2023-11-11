package commands

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/xopoww/standup/internal/bot/commands/commandtypes"
	"github.com/xopoww/standup/internal/bot/formatting"
	"github.com/xopoww/standup/pkg/config"
)

//go:embed commands.yaml
var commandsData string

func LoadDescriptions() ([]commandtypes.Desc, error) {
	return loadFrom(strings.NewReader(commandsData))
}

func loadFrom(r io.Reader) ([]commandtypes.Desc, error) {
	var cmds struct {
		Commands []commandtypes.Desc `yaml:"commands" validate:"required"`
	}
	err := config.Load(r, &cmds)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	for i := range cmds.Commands {
		var err error
		cmds.Commands[i].Short, err = formatting.RenderTemplate(cmds.Commands[i].Short, nil)
		if err != nil {
			return nil, fmt.Errorf("render commands[%d].short: %w", i, err)
		}
		cmds.Commands[i].Long, err = formatting.RenderTemplate(cmds.Commands[i].Long, nil)
		if err != nil {
			return nil, fmt.Errorf("render commands[%d].long: %w", i, err)
		}

		cmds.Commands[i].Usage = strings.TrimSpace(cmds.Commands[i].Usage)
		cmds.Commands[i].Short = strings.TrimSpace(cmds.Commands[i].Short)
		cmds.Commands[i].Long = strings.TrimSpace(cmds.Commands[i].Long)
	}
	return cmds.Commands, nil
}
