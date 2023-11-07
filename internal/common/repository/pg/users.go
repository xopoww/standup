package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/xopoww/standup/internal/bot/models"
	"github.com/xopoww/standup/internal/common/repository/dberrors"
)

var _ models.Models = &Repository{}

func (r *Repository) GetUser(ctx context.Context, username string) (*models.User, error) {
	const stmt = "get_user"
	_, err := r.conn.Prepare(ctx, stmt, "SELECT whitelisted FROM users WHERE username = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare: %w", err)
	}
	row := r.conn.QueryRow(ctx, stmt, username)
	user := &models.User{Username: username}
	err = row.Scan(&user.Whitelisted)
	if errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("user %q %w", username, dberrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpsertUser(ctx context.Context, user *models.User) error {
	const stmt = "upsert_user"
	_, err := r.conn.Prepare(ctx, stmt, "INSERT INTO users (username, whitelisted) VALUES  ($1, $2)"+
		"ON CONFLICT (username) DO UPDATE SET whitelisted=$2")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	_, err = r.conn.Exec(ctx, stmt, user.Username, user.Whitelisted)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
