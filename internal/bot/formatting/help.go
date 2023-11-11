package formatting

import (
	"github.com/xopoww/standup/internal/bot/commands/commandtypes"
)

func FormatHelp(cmds []commandtypes.Desc) string {
	data := struct {
		Text string
		Cmds []commandtypes.Desc

		FormatCommandShortHelp func(cmd commandtypes.Desc) string
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

func FormatCommandUsage(cmd commandtypes.Desc) string {
	usage := "/" + cmd.Name
	if cmd.Usage != "" {
		usage += " " + cmd.Usage
	}
	return Mono(usage)
}

func FormatCommandShortHelp(cmd commandtypes.Desc) string {
	return FormatCommandUsage(cmd) + " " + cmd.Short
}

func FormatCommandHelp(cmd commandtypes.Desc) string {
	help := FormatCommandShortHelp(cmd)
	if cmd.Long != "" {
		help += "\n\n" + cmd.Long
	}
	return help
}
