package identifiers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/pkg/identifiers"
)

func TestGenerate(t *testing.T) {
	id, err := identifiers.GenerateID()
	require.NoError(t, err)
	require.Len(t, id, identifiers.IDLength)

	id, err = identifiers.GenerateID()
	require.NoError(t, err)
	require.Len(t, id, identifiers.ShortIDLength)
}
