package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, dbs string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, dbs)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &Repository{conn: conn}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close(ctx)
}
