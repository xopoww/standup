package formatting

import (
	"github.com/xopoww/standup/pkg/api/standup"
)

// FormatMessages formats a list of messages for the report.
// If present, header must be valid when parsed with ParseMode (see Escape).
func FormatMessages(header string, msgs []*standup.Message) string {
	data := struct {
		Header   string
		Messages []*standup.Message
	}{
		Header:   header,
		Messages: msgs,
	}
	return MustRenderTemplate(`
{{- if .Header }}{{ bold .Header }}
{{ end }}
{{- if and .Header .Messages }}
{{ end }}{{ range .Messages }}- {{ esc .GetText }}
{{ end }}`, data)
}

func FormatMessageCreated(id string) string {
	return MustRenderTemplate(`Created message {{ mono . }}.`, id)
}
