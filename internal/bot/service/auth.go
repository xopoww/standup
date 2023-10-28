package service

import (
	"context"
	"time"

	"github.com/xopoww/standup/internal/grpcserver"
	"github.com/xopoww/standup/internal/logging"
	"google.golang.org/grpc/metadata"
)

const (
	ShortTTl = time.Minute * 5
)

func (s *Service) issueToken(ctx context.Context, user string, ttl time.Duration) (context.Context, error) {
	nb := time.Now().UTC()
	eat := nb.Add(ttl)
	token, err := s.deps.Issuer.IssueToken(user, nb, eat)
	if err != nil {
		return ctx, err
	}
	logging.L(ctx).Sugar().Debugf("Issued token for user %q valid until %s.", user, eat)
	return metadata.AppendToOutgoingContext(ctx, grpcserver.MetadataTokenKey, token), nil
}
