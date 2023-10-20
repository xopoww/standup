package grpcserver

import (
	"context"
	"time"

	"github.com/xopoww/standup/internal/models"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/identifiers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	standup.UnimplementedStandupServer

	repo models.Repository
}

func NewService(repo models.Repository) standup.StandupServer {
	return &service{repo: repo}
}

func (s *service) CreateMessage(ctx context.Context, req *standup.CreateMessageRequest) (*standup.CreateMessageResponse, error) {
	id := identifiers.GenerateID()
	msg := &models.Message{
		ID:        id,
		OwnerID:   req.OwnerId,
		Text:      req.Text,
		CreatedAt: time.Now(),
	}
	err := s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &standup.CreateMessageResponse{Id: id}, nil
}

func (s *service) GetMessage(ctx context.Context, req *standup.GetMessageRequest) (*standup.GetMessageResponse, error) {
	msg, err := s.repo.GetMessage(ctx, req.Id)
	if err != nil {
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
	msgs, err := s.repo.ListMessages(ctx, req.OwnerId, req.From.AsTime(), req.To.AsTime())
	if err != nil {
		return nil, err
	}
	resp := &standup.ListMessagesResponse{}
	for _, msg := range msgs {
		resp.Messages = append(resp.Messages, &standup.Message{
			Id:        msg.ID,
			Text:      msg.Text,
			OwnerId:   msg.OwnerID,
			CreatedAt: timestamppb.New(msg.CreatedAt),
		})
	}
	return resp, nil
}
