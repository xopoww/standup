package commands_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/commands"
)

func TestDescriptions(t *testing.T) {
	_, err := commands.LoadDescriptions()
	require.NoError(t, err)
}
