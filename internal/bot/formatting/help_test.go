package formatting_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/commands"
	"github.com/xopoww/standup/internal/bot/formatting"
)

func TestFormatCommandShortHelp(t *testing.T) {
	cc := []struct {
		name string
		cmd  commands.Desc
		want string
	}{
		{
			name: "simple",
			cmd: commands.Desc{
				Name:  "foo",
				Short: "test desc",
			},
			want: "`/foo` test desc",
		},
		{
			name: "with_usage",
			cmd: commands.Desc{
				Name:  "foo",
				Short: "test desc",
				Usage: "<arg1> [arg2]",
			},
			want: "`/foo <arg1> [arg2]` test desc",
		},
		{
			name: "escape",
			cmd: commands.Desc{
				Name:  "foo",
				Short: "test desc with dot.",
			},
			want: "`/foo` test desc with dot\\.",
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			s := formatting.FormatCommandShortHelp(c.cmd)
			require.Equal(t, c.want, s)
		})
	}
}

func TestFormatHelp(t *testing.T) {
	cmds := []commands.Desc{
		{
			Name:  "foo",
			Short: "foo command",
		},
		{
			Name:  "bar",
			Short: "bar command",
			Usage: "<arg1>",
		},
		{
			Name:  "baz",
			Short: "baz command",
			Long:  "A long description of bar command.\nTakes up several lines.",
		},
	}

	s := formatting.FormatHelp(cmds)
	require.Contains(t, s, "This bot can be used")
	require.Contains(t, s, "Availible commands \\(run `/help <command>` for more info\\):\n\n")
	require.Contains(t, s, "\\- `/foo` foo command\n")
	require.Contains(t, s, "\\- `/bar <arg1>` bar command\n")
	require.Contains(t, s, "\\- `/baz` baz command\n")
}
