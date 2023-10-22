package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/golang-migrate/migrate/v4"
	"github.com/xopoww/standup/internal/grpcserver"
	"github.com/xopoww/standup/internal/models"
	"github.com/xopoww/standup/internal/repository/pg"
	"github.com/xopoww/standup/internal/repository/pg/migrations"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc"
)

type Daemon struct {
	cfg Config

	repo models.Repository
	srv  *grpc.Server
}

func NewDaemon(ctx context.Context, cfg Config) (*Daemon, error) {
	mig, err := migrations.NewMigration(ctx, cfg.Database.DBS)
	if err != nil {
		return nil, fmt.Errorf("new migration: %w", err)
	}
	err = mig.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	repo, err := pg.NewRepository(ctx, cfg.Database.DBS)
	if err != nil {
		return nil, fmt.Errorf("new repo: %w", err)
	}

	grpcServer := grpc.NewServer()
	standup.RegisterStandupServer(grpcServer, grpcserver.NewService(repo))

	return &Daemon{
		cfg:  cfg,
		repo: repo,
		srv:  grpcServer,
	}, nil
}

func (d *Daemon) Start() error {
	lis, err := net.Listen("tcp", d.cfg.Service.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return d.srv.Serve(lis)
}

func (d *Daemon) GracefulStop(ctx context.Context) error {
	d.srv.GracefulStop()
	return d.repo.Close(ctx)
}
