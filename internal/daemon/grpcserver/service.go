//nolint:lll // grpc signatures
package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xopoww/standup/internal/common/auth"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/internal/common/repository/dberrors"
	"github.com/xopoww/standup/internal/daemon/models"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/identifiers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	standup.UnimplementedStandupServer

	repo models.Models

	ath auth.Authenticator
}

func NewService(repo models.Models, ath auth.Authenticator) standup.StandupServer {
	return &service{repo: repo, ath: ath}
}

func (s *service) CreateMessage(ctx context.Context, req *standup.CreateMessageRequest) (*standup.CreateMessageResponse, error) {
	if err := s.authorize(ctx, req.GetOwnerId()); err != nil {
		return nil, err
	}
	id, err := identifiers.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("generate id: %w", err)
	}
	msg := &models.Message{
		ID:        id,
		OwnerID:   req.GetOwnerId(),
		Text:      req.GetText(),
		CreatedAt: time.Now().UTC(),
	}
	err = s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, s.mapError(err)
	}
	logging.L(ctx).Sugar().Debugf("Created new message %q at %s.", id, msg.CreatedAt)
	return &standup.CreateMessageResponse{Id: id}, nil
}

func (s *service) GetMessage(ctx context.Context, req *standup.GetMessageRequest) (*standup.GetMessageResponse, error) {
	msg, err := s.repo.GetMessage(ctx, req.GetId())
	if err != nil {
		return nil, s.mapError(err)
	}
	if err := s.authorize(ctx, msg.OwnerID); err != nil {
		return nil, err
	}
	return &standup.GetMessageResponse{
		Message: &standup.Message{
			Id:        msg.ID,
			Text:      msg.Text,
			OwnerId:   msg.OwnerID,
			CreatedAt: timestamppb.New(msg.CreatedAt),
		},
	}, nil
}

func (s *service) ListMessages(ctx context.Context, req *standup.ListMessagesRequest) (*standup.ListMessagesResponse, error) {
	if err := s.authorize(ctx, req.GetOwnerId()); err != nil {
		return nil, err
	}
	msgs, err := s.repo.ListMessages(ctx, req.GetOwnerId(), req.GetFrom().AsTime(), req.GetTo().AsTime())
	if err != nil {
		return nil, s.mapError(err)
	}
	resp := &standup.ListMessagesResponse{}
	for _, msg := range msgs {
		resp.Messages = append(resp.GetMessages(), &standup.Message{
			Id:        msg.ID,
			Text:      msg.Text,
			OwnerId:   msg.OwnerID,
			CreatedAt: timestamppb.New(msg.CreatedAt),
		})
	}
	return resp, nil
}

func (s *service) mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dberrors.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
