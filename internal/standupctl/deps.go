package standupctl

import (
	"context"

	"github.com/xopoww/standup/internal/common/repository/pg"
)

type Deps struct {
	Cfg  Config
	Repo *pg.Repository
}

func NewDeps(ctx context.Context, cfg Config) (Deps, error) {
	repo, err := pg.NewRepository(ctx, cfg.Database.DBS)
	if err != nil {
		return Deps{}, err
	}
	return Deps{Cfg: cfg, Repo: repo}, nil
}
