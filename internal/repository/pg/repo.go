package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/xopoww/standup/internal/models"
)

type repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, dbs string) (models.Repository, error) {
	conn, err := pgx.Connect(ctx, dbs)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &repository{conn: conn}, nil
}

func (r *repository) Close(ctx context.Context) error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close(ctx)
}
