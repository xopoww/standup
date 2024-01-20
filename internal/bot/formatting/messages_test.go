package formatting_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/formatting"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFormatMessages(t *testing.T) {
	cc := []struct {
		name     string
		header   string
		messages []*standup.Message
		want     string
	}{
		{
			name:     "simple",
			messages: protoMessages("a msg"),
			want:     "\\- a msg\n",
		},
		{
			name:     "several",
			messages: protoMessages("a msg 1", "a msg 2", "a msg 3"),
			want:     "\\- a msg 1\n\\- a msg 2\n\\- a msg 3\n",
		},
		{
			name:     "header",
			header:   "A header",
			messages: protoMessages("a msg"),
			want:     "**A header**\n\n\\- a msg\n",
		},
		{
			name:     "escape_msg",
			messages: protoMessages("msg with dot."),
			want:     "\\- msg with dot\\.\n",
		},
		{
			name:   "only_header",
			header: "A header",
			want:   "**A header**\n",
		},
		{
			name: "empty",
			want: "",
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			s := formatting.FormatMessages(c.header, c.messages)
			require.Equal(t, c.want, s)
		})
	}
}

func TestFormatMessageCreated(t *testing.T) {
	s := formatting.FormatMessageCreated("testid")
	require.Equal(t, "Created message `testid`\\.", s)
}

func protoMessages(texts ...string) []*standup.Message {
	msgs := make([]*standup.Message, len(texts))
	for i := range texts {
		msgs[i] = &standup.Message{
			Id:        "test-id",
			Text:      texts[i],
			OwnerId:   12345,
			CreatedAt: timestamppb.Now(),
		}
	}
	return msgs
}
