package migrate

import "github.com/spf13/cobra"

func Migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
	}
	cmd.AddCommand(up())
	return cmd
}
