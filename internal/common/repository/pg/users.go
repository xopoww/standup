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

func (r *Repository) SetUserID(ctx context.Context, username string, id int64) error {
	ct, err := r.conn.Exec(ctx, `UPDATE users SET id = $1 WHERE username = $2 AND id = NULL`, id, username)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return dberrors.ErrNotFound
	}
	return nil
}

func (r *Repository) UpsertUser(ctx context.Context, user *models.User) error {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO users (id, username, whitelisted) VALUES ($1, $2, false)
		ON CONFLICT (id) DO UPDATE SET username = $2 WHERE id = $1
	`, user.ID, user.Username)
	return err
}

func (r *Repository) GetWhitelisted(ctx context.Context, userID int64) (bool, error) {
	row := r.conn.QueryRow(ctx, `SELECT whitelisted FROM users WHERE id = $1`, userID)
	whitelisted := false
	err := row.Scan(&whitelisted)
	if errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("user %d %w", userID, dberrors.ErrNotFound)
	}
	return whitelisted, err
}

func (r *Repository) SetWhitelisted(ctx context.Context, userID int64, whitelisted bool) error {
	ct, err := r.conn.Exec(ctx, `UPDATE users SET whitelisted = $1 WHERE id = $2`, whitelisted, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user %d %w", userID, dberrors.ErrNotFound)
	}
	return nil
}
