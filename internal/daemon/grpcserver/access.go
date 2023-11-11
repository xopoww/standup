package grpcserver

import (
	"context"

	"github.com/xopoww/standup/internal/common/auth"
	"github.com/xopoww/standup/internal/common/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *service) authenticate(ctx context.Context) (userID string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}
	mdToken, ok := md[auth.GRPCMetadataTokenKey]
	if !ok || len(mdToken) < 1 {
		return "", status.Errorf(codes.Unauthenticated, "missing %q metadata field", auth.GRPCMetadataTokenKey)
	}
	userID, err = s.ath.Authenticate(mdToken[0])
	if err != nil {
		logging.L(ctx).Sugar().Warnf("Authentication failed: %s.", err)
		return "", status.Errorf(codes.Unauthenticated, "bad token")
	}
	return userID, nil
}

func (s *service) authorize(ctx context.Context, targetUserID string) error {
	if !s.ath.Enabled() {
		return nil
	}
	userID, err := s.authenticate(ctx)
	if err != nil {
		return err
	}
	if userID != targetUserID {
		return status.Errorf(codes.PermissionDenied, "permission denied for user %q", userID)
	}
	return nil
}
