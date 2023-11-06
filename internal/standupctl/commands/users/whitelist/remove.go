package whitelist

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/bot/models"
	"github.com/xopoww/standup/internal/standupctl"
)

func remove(deps *standupctl.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <username>",
		Aliases: []string{"rm"},
		Short:   "Remove user from whitelist",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, pos []string) error {
			username := pos[0]
			return deps.Repo.UpsertUser(context.Background(), &models.User{
				Username:    username,
				Whitelisted: false,
			})
		},
	}
	return cmd
}
