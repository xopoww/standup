package secrets

import "github.com/spf13/cobra"

func Secrets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage secrets",
	}
	cmd.AddCommand(genKeys())
	return cmd
}
