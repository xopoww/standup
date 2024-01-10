package migrate

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
)

func Migrate(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Perform database migrations",
	}
	cmd.AddCommand(up(deps))
	cmd.AddCommand(force(deps))
	cmd.AddCommand(to(deps))
	return cmd
}
