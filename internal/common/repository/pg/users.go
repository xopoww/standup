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

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.conn.QueryRow(ctx, `SELECT username FROM users WHERE id = $1`, id)
	user := &models.User{ID: id}
	err := row.Scan(&user.Username)
	if errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("user %d %w", id, dberrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpsertUser(ctx context.Context, user *models.User) error {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO users (id, username) VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE SET username = $2
	`, user.ID, user.Username)
	return err
}

func (r *Repository) GetWhitelisted(ctx context.Context, userID int64) (bool, error) {
	row := r.conn.QueryRow(ctx, `SELECT whitelisted FROM user_whitelist WHERE user_id = $1`, userID)
	whitelisted := false
	err := row.Scan(&whitelisted)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return whitelisted, err
}

func (r *Repository) SetWhitelisted(ctx context.Context, userID int64, whitelisted bool) error {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO user_whitelist (user_id, whitelisted) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET whitelisted = $2
	`, userID, whitelisted)
	return err
}
