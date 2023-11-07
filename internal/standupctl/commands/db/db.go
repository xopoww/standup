package db

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
	"github.com/xopoww/standup/internal/standupctl/commands/db/migrate"
)

func DB(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Manage database",
	}
	cmd.AddCommand(migrate.Migrate(deps))
	return cmd
}
