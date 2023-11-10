package formatting

import (
	"fmt"
	"strings"

	"github.com/xopoww/standup/internal/bot/commands"
)

func FormatHelp(text string, cmds []commands.Desc) string {
	bldr := &strings.Builder{}

	_, _ = fmt.Fprintf(bldr, "**Help**\n\n")
	if text != "" {
		_, _ = fmt.Fprintf(bldr, "%s\n\n", text)
	}
	_, _ = fmt.Fprintf(bldr, "")
	for _, cmd := range cmds {
		_, _ = fmt.Fprintf(bldr, "- %s\n", FormatCommandShortHelp(cmd))
	}
	return bldr.String()
}

func FormatCommandShortHelp(cmd commands.Desc) string {
	if cmd.Usage != "" {
		return fmt.Sprintf("`/%s %s` %s", cmd.Name, cmd.Usage, cmd.Short)
	}
	return fmt.Sprintf("`/%s` %s", cmd.Name, cmd.Short)
}

func FormatCommandHelp(cmd commands.Desc) string {
	help := FormatCommandShortHelp(cmd)
	if cmd.Long != "" {
		help += "\n\n" + cmd.Long
	}
	return help
}
