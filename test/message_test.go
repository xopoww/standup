package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/testutil"
	"github.com/xopoww/standup/pkg/api/standup"
)

func TestCreateMessage(t *testing.T) {
	ctx, cancel := testutil.NewContext(context.Background())
	defer cancel()
	t.Logf("TestID: %s.", testutil.TestID(ctx))

	req := &standup.CreateMessageRequest{
		Text:    "Test message",
		OwnerId: "test-owner-" + testutil.TestID(ctx),
	}
	rsp, err := deps.client.CreateMessage(ctx, req)
	require.NoError(t, err)

	var content, owner string
	row := deps.db.QueryRow(ctx, "SELECT content, owner_id FROM messages WHERE id = $1", rsp.Id)
	require.NoError(t, row.Scan(&content, &owner))
	require.Equal(t, req.Text, content)
	require.Equal(t, req.OwnerId, owner)
}
