package db

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl/commands/db/migrate"
)

func DB() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Manage database",
	}
	cmd.AddCommand(migrate.Migrate())
	return cmd
}
