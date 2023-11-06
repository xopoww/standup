package users

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
	"github.com/xopoww/standup/internal/standupctl/commands/users/whitelist"
)

func Users(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users",
		Aliases: []string{"user"},
		Short:   "Manage users",
	}
	cmd.AddCommand(whitelist.Whitelist(deps))
	return cmd
}
