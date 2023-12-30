package whitelist

import (
	"context"
	"fmt"
	"strconv"

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

func run(deps *standupctl.Deps, add bool) func(_ *cobra.Command, pos []string) error {
	return func(_ *cobra.Command, pos []string) error {
		id, err := strconv.ParseInt(pos[0], 10, 64)
		if err != nil {
			return fmt.Errorf("ParseInt: %w", err)
		}
		return deps.Repo.SetWhitelisted(context.Background(), id, add)
	}
}
