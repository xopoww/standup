package models

import "context"

type Models interface {
	GetUser(ctx context.Context, username string) (*User, error)
}
