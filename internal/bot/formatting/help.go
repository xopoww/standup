package formatting

import (
	"github.com/xopoww/standup/internal/bot/commands"
)

func FormatHelp(cmds []commands.Desc) string {
	data := struct {
		Text string
		Cmds []commands.Desc

		FormatCommandShortHelp func(cmd commands.Desc) string
	}{
		Cmds:                   cmds,
		FormatCommandShortHelp: FormatCommandShortHelp,
	}
	return MustRenderTemplate(`This bot can be used to save short messages and retrieve time-based reports.

Send a message to this bot to save it.

Availible commands (run {{ mono "/help <command>" }} for more info):

{{ range .Cmds }}- {{ call $.FormatCommandShortHelp . }}
{{end}}`, data)
}

func FormatCommandShortHelp(cmd commands.Desc) string {
	data := struct {
		Command string
		Short   string
	}{
		Command: "/" + cmd.Name,
		Short:   cmd.Short,
	}
	if cmd.Usage != "" {
		data.Command += " " + cmd.Usage
	}
	return MustRenderTemplate(`{{ mono .Command }} {{ esc .Short }}`, data)
}

func FormatCommandHelp(cmd commands.Desc) string {
	help := FormatCommandShortHelp(cmd)
	if cmd.Long != "" {
		help += "\n\n" + Escape(cmd.Long)
	}
	return help
}
