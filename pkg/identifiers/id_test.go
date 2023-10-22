package identifiers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/pkg/identifiers"
)

func TestGenerate(t *testing.T) {
	require.Equal(t, 16, len(identifiers.GenerateID()))
	require.Equal(t, 8, len(identifiers.GenerateShortID()))
}
