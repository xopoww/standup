package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrations embed.FS

func NewMigration(ctx context.Context, dbs string) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations, ".")
	if err != nil {
		return nil, fmt.Errorf("new source: %w", err)
	}
	defer src.Close()
	return migrate.NewWithSourceInstance("iofs", src, dbs)
}
