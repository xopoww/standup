package whitelist

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
)

func add(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <id>",
		Short: "Add user to whitelist",
		Args:  cobra.ExactArgs(1),
		RunE:  run(deps, true),
	}
	return cmd
}
