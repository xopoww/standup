package formatting

import (
	"fmt"
	"strings"

	"github.com/xopoww/standup/pkg/api/standup"
)

func FormatMessages(header string, msgs []*standup.Message) string {
	bldr := &strings.Builder{}
	if header != "" {
		_, _ = fmt.Fprintf(bldr, "**%s**\n\n", header)
	}
	for _, msg := range msgs {
		_, _ = fmt.Fprintf(bldr, "- %s\n", msg.GetText())
	}
	return bldr.String()
}
