package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/common/testutil"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/identifiers"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateMessage(t *testing.T) {
	RunTest(t, "default", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))

		req := testutil.CreateMessageRequest(ctx)
		rsp, err := deps.Client.CreateMessage(ctx, req)
		require.NoError(t, err)

		var content, owner string
		row := deps.DB.QueryRow(ctx, "SELECT content, owner_id FROM messages WHERE id = $1", rsp.GetId())
		require.NoError(t, row.Scan(&content, &owner))
		require.Equal(t, req.GetText(), content)
		require.Equal(t, req.GetOwnerId(), owner)
	})
}

func TestGetMessage(t *testing.T) {
	RunTest(t, "default", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))

		req := testutil.CreateMessageRequest(ctx)
		rsp, err := deps.Client.CreateMessage(ctx, req)
		require.NoError(t, err)

		msg, err := deps.Client.GetMessage(ctx, &standup.GetMessageRequest{Id: rsp.GetId()})
		require.NoError(t, err)
		require.Equal(t, req.GetText(), msg.GetMessage().GetText())
		require.Equal(t, req.GetOwnerId(), msg.GetMessage().GetOwnerId())
	})

	RunTest(t, "not_found", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))
		id, err := identifiers.GenerateID()
		require.NoError(t, err)
		_, err = deps.Client.GetMessage(ctx, &standup.GetMessageRequest{Id: id})
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
			rsp, err := deps.Client.CreateMessage(ctx, req)
			require.NoError(t, err)
			expected[i] = rsp.GetId()
		}

		rsp, err := deps.Client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(from),
			To:      timestamppb.Now(),
		})
		require.NoError(t, err)
		var actual []string
		for _, msg := range rsp.GetMessages() {
			actual = append(actual, msg.GetId())
		}
		require.Equal(t, expected, actual)
	})

	RunTest(t, "partial", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))
		from := time.Now()

		req := testutil.CreateMessageRequest(ctx)
		msg, err := deps.Client.CreateMessage(ctx, req)
		require.NoError(t, err)

		to := time.Now()

		_, err = deps.Client.CreateMessage(ctx, req)
		require.NoError(t, err)

		rsp, err := deps.Client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(from),
			To:      timestamppb.New(to),
		})
		require.NoError(t, err)
		require.Len(t, rsp.GetMessages(), 1)
		require.Equal(t, msg.GetId(), rsp.GetMessages()[0].GetId())
	})

	RunTest(t, "another_owner", func(ctx context.Context, t *testing.T) {
		from := time.Now()

		req := testutil.CreateMessageRequest(ctx)
		msg, err := deps.Client.CreateMessage(withToken(ctx, t, req.GetOwnerId()), req)
		require.NoError(t, err)

		req.OwnerId = "another-" + req.GetOwnerId()
		_, err = deps.Client.CreateMessage(withToken(ctx, t, req.GetOwnerId()), req)
		require.NoError(t, err)

		rsp, err := deps.Client.ListMessages(withToken(ctx, t, testutil.OwnerID(ctx)), &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(from),
			To:      timestamppb.Now(),
		})
		require.NoError(t, err)
		require.Len(t, rsp.GetMessages(), 1)
		require.Equal(t, msg.GetId(), rsp.GetMessages()[0].GetId())
	})

	RunTest(t, "empty", func(ctx context.Context, t *testing.T) {
		ctx = withToken(ctx, t, testutil.OwnerID(ctx))
		rsp, err := deps.Client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(time.Now()),
			To:      timestamppb.New(time.Now().Add(time.Hour)),
		})
		require.NoError(t, err)
		require.Empty(t, rsp.GetMessages())
	})
}
