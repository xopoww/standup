package service

import (
	"context"
	"time"

	"github.com/xopoww/standup/internal/common/auth"
	"github.com/xopoww/standup/internal/common/logging"
	"google.golang.org/grpc/metadata"
)

const (
	ShortTTL = time.Minute * 5
)

func (s *Service) issueToken(ctx context.Context, userID int64, ttl time.Duration) (context.Context, error) {
	nb := time.Now().UTC()
	eat := nb.Add(ttl)
	token, err := s.deps.Issuer.IssueToken(userID, nb, eat)
	if err != nil {
		return ctx, err
	}
	logging.L(ctx).Sugar().Debugf("Issued token for user %q valid until %s.", userID, eat)
	return metadata.AppendToOutgoingContext(ctx, auth.GRPCMetadataTokenKey, token), nil
}
