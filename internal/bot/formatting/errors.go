package formatting

import (
	"fmt"

	"github.com/xopoww/standup/internal/bot/commands/commandtypes"
)

type Printable interface {
	Printable() string
}

type SyntaxError struct {
	Msg string
	Cmd *commandtypes.Desc
}

func NewSyntaxErrorf(format string, a ...any) SyntaxError {
	return SyntaxError{Msg: fmt.Sprintf(format, a...)}
}

func NewCommandSyntaxErrorf(cmd commandtypes.Desc, format string, a ...any) SyntaxError {
	e := NewSyntaxErrorf(format, a...)
	e.Cmd = &cmd
	return e
}

func (e SyntaxError) Error() string {
	s := e.Msg
	if e.Cmd != nil {
		s = fmt.Sprintf("%s: %s", e.Cmd.Name, s)
	}
	return s
}

func (e SyntaxError) Printable() string {
	data := map[string]any{
		"E": e,
		"FormatCommandUsage": func(c *commandtypes.Desc) string {
			return FormatCommandUsage(*c)
		},
	}
	return MustRenderTemplate(`Error: {{ esc .E.Msg }}.{{ if .E.Cmd }}

Usage: {{ call .FormatCommandUsage .E.Cmd }}.{{end}}`, data)
}

var (
	InternalError    = Escape("Internal error occured.")
	PermissionDenied = Escape("Permission denied.")
	UsageRestricted  = Escape("Bot usage is currently restricted.")
)
