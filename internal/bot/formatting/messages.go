package formatting

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/pkg/api/standup"
)

const ParseMode = tgbotapi.ModeMarkdownV2

// Escape makes arbitrary text valid when parsed with ParseMode.
func Escape(s string) string {
	return tgbotapi.EscapeText(ParseMode, s)
}

// FormatMessages formats a list of messages for the report.
// If present, header must be valid when parsed with ParseMode (see Escape).
func FormatMessages(header string, msgs []*standup.Message) string {
	bldr := &strings.Builder{}
	if header != "" {
		_, _ = fmt.Fprintf(bldr, "*%s*\n\n", header)
	}
	for _, msg := range msgs {
		_, _ = fmt.Fprintf(bldr, "\\- %s\n", Escape(msg.GetText()))
	}
	return bldr.String()
}
