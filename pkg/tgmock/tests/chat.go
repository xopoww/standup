package tests

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/pkg/tgmock/control"
)

func ChatHistory(ctx context.Context, t require.TestingT, tm control.TGMockControlClient, c *control.Chat) string {
	rsp, err := tm.ListMessages(ctx, &control.ListMessagesRequest{
		ChatId: c.GetId(),
		All:    true,
	})
	require.NoError(t, err)

	bldr := &strings.Builder{}
	for _, msg := range rsp.GetMessages() {
		date := time.Unix(msg.GetDate(), 0).Local()
		fmt.Fprintf(bldr, headerFormat, date.Format(dateFormat), msg.GetFrom().GetUsername())
		if msg.GetReplyToMessage() != nil {
			writePaddedText(bldr, "> %s\n", renderMarkdownV2Text(msg.GetReplyToMessage().GetText()))
		}
		writePaddedText(bldr, "%s\n", renderMarkdownV2Text(msg.GetText()))
	}
	return bldr.String()
}

func renderMarkdownV2Text(text string) string {
	bldr := &strings.Builder{}

	escaped := false
	for _, b := range text {
		if escaped {
			bldr.WriteRune(b)
			escaped = false
			continue
		}
		switch b {
		case '\\':
			escaped = true
			continue
		default:
			bldr.WriteRune(b)
		}
	}
	return bldr.String()
}

func writePaddedText(w io.Writer, lineFormat, text string) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		fmt.Fprintf(w, lineFormat, line)
		if i+1 < len(lines) {
			fmt.Fprint(w, padding)
		}
	}
}

const (
	basePadding   = "\t\t"
	dateFormat    = "02-01-06 15:04:05"
	dateWidth     = len(dateFormat) + 1
	usernameWidth = 32 + 1
)

var (
	headerFormat = basePadding + fmt.Sprintf("%%%ds", dateWidth) + fmt.Sprintf("%%%ds", usernameWidth) + " : "
	padding      = basePadding + strings.Repeat(" ", dateWidth+usernameWidth) + "   "
)
