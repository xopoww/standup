package models

import "context"

type Models interface {
	GetUserByID(ctx context.Context, id int64) (*User, error)

	// TODO: rm after transition period
	SetUserID(ctx context.Context, username string, id int64) error

	UpsertUser(ctx context.Context, user *User) error

	GetWhitelisted(ctx context.Context, userID int64) (bool, error)
	SetWhitelisted(ctx context.Context, userID int64, whitelisted bool) error
}
