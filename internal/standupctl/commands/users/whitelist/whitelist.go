package whitelist

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
)

func Whitelist(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist",
		Short: "Manage user whitelist",
	}
	cmd.AddCommand(add(deps), remove(deps))
	return cmd
}
