package db

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl/commands/db/migrate"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "db",
	}
	cmd.AddCommand(migrate.Migrate())
	return cmd
}
