package service

import "context"

type Repository interface {
	CheckWhitelist(ctx context.Context, username string) (bool, error)
}
