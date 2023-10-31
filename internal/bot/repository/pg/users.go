package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) CheckWhitelist(ctx context.Context, username string) (bool, error) {
	const stmt = "get_message"
	_, err := r.conn.Prepare(ctx, stmt, "SELECT allowed FROM users WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("prepare: %w", err)
	}
	row := r.conn.QueryRow(ctx, stmt, username)
	var allowed bool
	err = row.Scan(&allowed)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return allowed, nil
}
