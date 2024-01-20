package models

import (
	"context"
	"time"
)

type Models interface {
	CreateMessage(ctx context.Context, msg *Message) error
	GetMessage(ctx context.Context, id string) (*Message, error)
	ListMessages(ctx context.Context, ownerID int64, from, to time.Time) ([]*Message, error)
}
