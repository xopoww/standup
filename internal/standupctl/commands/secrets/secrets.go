package secrets

import (
	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
)

func Secrets(_ *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage secrets",
	}
	cmd.AddCommand(genKeys())
	return cmd
}
