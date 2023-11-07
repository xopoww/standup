package migrate

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/common/repository/pg/migrations"
	"github.com/xopoww/standup/internal/standupctl"
)

func up(deps *standupctl.Deps) *cobra.Command {
	var args struct {
		verbose bool
	}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate up",
		Long:  "Migrate the database to the most resent version",
		RunE: func(_ *cobra.Command, _ []string) error {
			mig, err := migrations.NewMigration(context.Background(), deps.Cfg.Database.DBS)
			if err != nil {
				return fmt.Errorf("new migration: %w", err)
			}
			mig.Log = migrateLogger{args.verbose}

			err = mig.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				err = nil
			}
			return err
		},
	}
	cmd.Flags().BoolVarP(&args.verbose, "verbose", "v", false, "")
	return cmd
}

type migrateLogger struct {
	verbose bool
}

func (l migrateLogger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func (l migrateLogger) Verbose() bool {
	return l.verbose
}
