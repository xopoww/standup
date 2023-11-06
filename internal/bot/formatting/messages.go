package formatting

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xopoww/standup/pkg/api/standup"
)

const ParseMode = tgbotapi.ModeMarkdownV2

func FormatMessages(header string, msgs []*standup.Message) string {
	bldr := &strings.Builder{}
	if header != "" {
		_, _ = fmt.Fprintf(bldr, "*%s*\n\n", header)
	}
	for _, msg := range msgs {
		_, _ = fmt.Fprintf(bldr, "\\- %s\n", tgbotapi.EscapeText(ParseMode, msg.GetText()))
	}
	return bldr.String()
}
