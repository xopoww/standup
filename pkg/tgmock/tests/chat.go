package tests

import (
	"context"
	"fmt"
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
		fmt.Fprintf(bldr, "\t\t%17s %32s : ", date.Format("02-01-06 15:04:05"), msg.GetFrom().GetUsername())
		if msg.GetReplyToMessage() != nil {
			fmt.Fprintf(bldr, "> %s\n", msg.GetReplyToMessage().GetText())
			fmt.Fprintf(bldr, "\t\t%17s %32s   ", "", "")
		}
		fmt.Fprintf(bldr, "%s\n", msg.GetText())
	}
	return bldr.String()
}
