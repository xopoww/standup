package migrate

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/common/repository/pg/migrations"
	"github.com/xopoww/standup/internal/standupctl"
)

func force(deps *standupctl.Deps) *cobra.Command {
	var args struct {
		verbose bool
		version int
	}
	cmd := &cobra.Command{
		Use:   "force",
		Short: "Force database version",
		Long:  "Forcefully set database version (use with caution)",
		RunE: func(_ *cobra.Command, _ []string) error {
			mig, err := migrations.NewMigration(context.Background(), deps.Cfg.Database.DBS)
			if err != nil {
				return fmt.Errorf("new migration: %w", err)
			}
			mig.Log = migrateLogger{args.verbose}

			return mig.Force(args.version)
		},
	}
	cmd.Flags().BoolVarP(&args.verbose, "verbose", "v", false, "")
	cmd.Flags().IntVar(&args.version, "version", -1, "DB version")
	return cmd
}
