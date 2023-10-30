package commands

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl/commands/db"
)

func Root() *cobra.Command {
	root := &cobra.Command{
		Use: "standupctl",
	}
	root.AddCommand(db.Command())

	return root
}
