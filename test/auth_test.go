package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/testutil"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuth(t *testing.T) {
	if !deps.cfg.Auth.Enabled {
		t.Skip("authentication is disabled")
	}

	RunTest(t, "wrong_subject", func(ctx context.Context, t *testing.T) {
		msg, err := deps.client.CreateMessage(withToken(ctx, t, testutil.OwnerID(ctx)), testutil.CreateMessageRequest(ctx))
		require.NoError(t, err)

		ctx = withToken(ctx, t, "wrong-"+testutil.OwnerID(ctx))

		_, err = deps.client.CreateMessage(ctx, testutil.CreateMessageRequest(ctx))
		testutil.RequireErrCode(t, codes.PermissionDenied, err)

		_, err = deps.client.GetMessage(ctx, &standup.GetMessageRequest{Id: msg.GetId()})
		testutil.RequireErrCode(t, codes.PermissionDenied, err)

		_, err = deps.client.ListMessages(ctx, &standup.ListMessagesRequest{
			OwnerId: testutil.OwnerID(ctx),
			From:    timestamppb.New(time.Now()),
			To:      timestamppb.New(time.Now().Add(time.Hour)),
		})
		testutil.RequireErrCode(t, codes.PermissionDenied, err)
	})
}
