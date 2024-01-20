package migrate

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/common/repository/pg/migrations"
	"github.com/xopoww/standup/internal/standupctl"
)

func to(deps *standupctl.Deps) *cobra.Command {
	var args struct {
		verbose bool
	}
	cmd := &cobra.Command{
		Use:   "to",
		Short: "Migrate DB",
		Long:  "Migrate the database to specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, pos []string) error {
			version, err := strconv.ParseUint(pos[0], 10, 32)
			if err != nil {
				return err
			}

			mig, err := migrations.NewMigration(context.Background(), deps.Cfg.Database.DBS)
			if err != nil {
				return fmt.Errorf("new migration: %w", err)
			}
			mig.Log = migrateLogger{args.verbose}

			err = mig.Migrate(uint(version))
			if errors.Is(err, migrate.ErrNoChange) {
				err = nil
			}
			return err
		},
	}
	cmd.Flags().BoolVarP(&args.verbose, "verbose", "v", false, "")
	return cmd
}
