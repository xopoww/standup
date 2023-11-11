package formatting_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/commands/commandtypes"
	"github.com/xopoww/standup/internal/bot/formatting"
)

func TestSyntaxError_Printable(t *testing.T) {
	t.Run("escape", func(t *testing.T) {
		e := formatting.NewSyntaxErrorf("some syntax error with-dash")
		require.Equal(t, "Error: some syntax error with\\-dash\\.", e.Printable())
	})

	t.Run("command", func(t *testing.T) {
		c := commandtypes.Desc{
			Name:  "foo",
			Usage: "<arg>",
		}
		e := formatting.NewCommandSyntaxErrorf(c, "some syntax error")
		require.Equal(t,
			"Error: some syntax error\\.\n\nUsage: `/foo <arg>`\\.",
			e.Printable(),
		)
	})
}
