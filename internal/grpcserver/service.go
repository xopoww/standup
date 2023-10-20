package grpcserver

import (
	"context"
	"log"

	"github.com/xopoww/standup/pkg/api/standup"
)

type service struct {
	standup.UnimplementedStandupServer
}

func NewService() standup.StandupServer {
	return &service{}
}

func (s *service) CreateMessage(ctx context.Context, req *standup.CreateMessageRequest) (*standup.CreateMessageResponse, error) {
	log.Printf("CreateMessage %q", req.Text)
	return &standup.CreateMessageResponse{
		Id: "foo",
	}, nil
}
