package commands

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl/commands/db"
	"github.com/xopoww/standup/internal/standupctl/commands/secrets"
)

func Root() *cobra.Command {
	root := &cobra.Command{
		Use: "standupctl",
	}
	root.AddCommand(db.DB(), secrets.Secrets())

	return root
}
