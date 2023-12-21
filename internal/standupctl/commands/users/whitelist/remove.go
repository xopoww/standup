package whitelist

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
)

func remove(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <id>",
		Aliases: []string{"rm"},
		Short:   "Remove user from whitelist",
		Args:    cobra.ExactArgs(1),
		RunE:    run(deps, false),
	}
	return cmd
}
