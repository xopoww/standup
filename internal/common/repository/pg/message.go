package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xopoww/standup/internal/common/repository/dberrors"
	"github.com/xopoww/standup/internal/daemon/models"
)

var _ models.Models = &Repository{}

func (r *Repository) CreateMessage(ctx context.Context, msg *models.Message) error {
	const stmt = "create_message"
	_, err := r.conn.Prepare(ctx, stmt, "INSERT INTO messages(id, owner_id, content, created_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	_, err = r.conn.Exec(ctx, stmt, msg.ID, msg.OwnerID, msg.Text, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repository) GetMessage(ctx context.Context, id string) (*models.Message, error) {
	const stmt = "get_message"
	_, err := r.conn.Prepare(ctx, stmt, "SELECT owner_id, content, created_at FROM messages WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare: %w", err)
	}
	row := r.conn.QueryRow(ctx, stmt, id)
	msg := &models.Message{ID: id}
	err = row.Scan(&msg.OwnerID, &msg.Text, &msg.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("message %q %w", id, dberrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *Repository) ListMessages(ctx context.Context, ownerID int64, from, to time.Time) ([]*models.Message, error) {
	const stmt = "list_messages"
	_, err := r.conn.Prepare(ctx, stmt,
		"SELECT id, content, created_at FROM messages WHERE owner_id = $1 AND created_at >= $2 AND created_at < $3",
	)
	if err != nil {
		return nil, fmt.Errorf("prepare: %w", err)
	}

	rows, err := r.conn.Query(ctx, stmt, ownerID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	var result []*models.Message
	for rows.Next() {
		msg := &models.Message{OwnerID: ownerID}
		err := rows.Scan(&msg.ID, &msg.Text, &msg.CreatedAt)
		if err != nil {
			return result, fmt.Errorf("scan: %w", err)
		}
		result = append(result, msg)
	}
	return result, rows.Err()
}
