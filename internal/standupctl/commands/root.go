package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/standupctl"
	"github.com/xopoww/standup/internal/standupctl/commands/db"
	"github.com/xopoww/standup/internal/standupctl/commands/secrets"
	"github.com/xopoww/standup/internal/standupctl/commands/users"
	"github.com/xopoww/standup/pkg/config"
)

func Root() *cobra.Command {
	var args struct {
		cfgPath string
	}
	deps := &standupctl.Deps{}
	root := &cobra.Command{
		Use: "standupctl",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			var cfg standupctl.Config
			err := config.LoadFile(args.cfgPath, &cfg)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			*deps, err = standupctl.NewDeps(context.Background(), cfg)
			if err != nil {
				return fmt.Errorf("new deps: %w", err)
			}
			return nil
		},
	}
	root.AddCommand(db.DB(deps), secrets.Secrets(deps), users.Users(deps))
	root.PersistentFlags().StringVar(&args.cfgPath, "config", standupctl.DefaultConfigPath, "path to yaml config file")
	_ = root.MarkPersistentFlagFilename("config")
	return root
}
