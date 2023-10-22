package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/testutil"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/identifiers"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateMessage(t *testing.T) {
	RunTest(t, "default", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))

		req := testutil.CreateMessageRequest(ctx)
		rsp, err := deps.client.CreateMessage(ctx, req)
		require.NoError(t, err)

		var content, owner string
		row := deps.db.QueryRow(ctx, "SELECT content, owner_id FROM messages WHERE id = $1", rsp.Id)
		require.NoError(t, row.Scan(&content, &owner))
		require.Equal(t, req.Text, content)
		require.Equal(t, req.OwnerId, owner)
	})
}

func TestGetMessage(t *testing.T) {
	RunTest(t, "default", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))

		req := testutil.CreateMessageRequest(ctx)
		rsp, err := deps.client.CreateMessage(ctx, req)
		require.NoError(t, err)

		msg, err := deps.client.GetMessage(ctx, &standup.GetMessageRequest{Id: rsp.Id})
		require.NoError(t, err)
		require.Equal(t, req.Text, msg.Message.Text)
		require.Equal(t, req.OwnerId, msg.Message.OwnerId)
	})

	RunTest(t, "not_found", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))

		_, err := deps.client.GetMessage(ctx, &standup.GetMessageRequest{Id: identifiers.GenerateID()})
		testutil.RequireErrCode(t, codes.NotFound, err)
	})
}

func TestListMessages(t *testing.T) {
	RunTest(t, "default", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))
		from := time.Now()

		expected := make([]string, 3)
		for i := range expected {
			req := testutil.CreateMessageRequest(ctx)
			rsp, err := deps.client.CreateMessage(ctx, req)
			require.NoError(t, err)
			expected[i] = rsp.Id
		}

		rsp, err := deps.client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(from),
			To:      timestamppb.Now(),
		})
		require.NoError(t, err)
		var actual []string
		for _, msg := range rsp.Messages {
			actual = append(actual, msg.Id)
		}
		require.Equal(t, expected, actual)
	})

	RunTest(t, "partial", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))
		from := time.Now()

		req := testutil.CreateMessageRequest(ctx)
		msg, err := deps.client.CreateMessage(ctx, req)
		require.NoError(t, err)

		to := time.Now()

		_, err = deps.client.CreateMessage(ctx, req)
		require.NoError(t, err)

		rsp, err := deps.client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(from),
			To:      timestamppb.New(to),
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(rsp.Messages))
		require.Equal(t, msg.Id, rsp.Messages[0].Id)
	})

	RunTest(t, "another_owner", func(ctx context.Context, t *testing.T) {
		from := time.Now()

		req := testutil.CreateMessageRequest(ctx)
		msg, err := deps.client.CreateMessage(withToken(ctx, t, req.OwnerId), req)
		require.NoError(t, err)

		req.OwnerId = "another-" + req.OwnerId
		_, err = deps.client.CreateMessage(withToken(ctx, t, req.OwnerId), req)
		require.NoError(t, err)

		rsp, err := deps.client.ListMessages(withToken(ctx, t, testutil.OwnerID(ctx)), &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(from),
			To:      timestamppb.Now(),
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(rsp.Messages))
		require.Equal(t, msg.Id, rsp.Messages[0].Id)
	})

	RunTest(t, "empty", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))
		rsp, err := deps.client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(time.Now()),
			To:      timestamppb.New(time.Now().Add(time.Hour)),
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(rsp.Messages))
	})
}
