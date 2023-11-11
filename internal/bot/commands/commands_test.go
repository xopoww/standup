package commands_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/commands"
	"github.com/xopoww/standup/internal/bot/formatting"
)

func TestDescriptions(t *testing.T) {
	cmds, err := commands.LoadDescriptions()
	require.NoError(t, err)

	t.Log(formatting.FormatHelp(cmds))
	for _, cmd := range cmds {
		t.Log(formatting.FormatCommandHelp(cmd))
	}
}
