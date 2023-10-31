package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireErrCode(t *testing.T, expected codes.Code, err error) {
	require.Error(t, err)
	s := status.Convert(err)
	require.Equalf(t, expected, s.Code(), "want %q, got %q", expected, s.Code())
}
