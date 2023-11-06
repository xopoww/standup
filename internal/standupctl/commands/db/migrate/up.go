package migrate

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/common/repository/pg/migrations"
)

func up() *cobra.Command {
	var args struct {
		dbs     string
		verbose bool
	}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate up",
		Long:  "Migrate the database to the most resent version",
		RunE: func(_ *cobra.Command, _ []string) error {
			mig, err := migrations.NewMigration(context.Background(), args.dbs)
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
	cmd.Flags().StringVar(&args.dbs, "dbs", "", "database connection string")
	cmd.Flags().BoolVarP(&args.verbose, "verbose", "v", false, "")
	_ = cmd.MarkFlagRequired("dbs")
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
