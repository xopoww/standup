package daemon

import (
	"context"
	"fmt"
	"net"

	"github.com/xopoww/standup/internal/common/auth"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/internal/common/repository/pg"
	"github.com/xopoww/standup/internal/daemon/grpcserver"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc"
)

type Daemon struct {
	cfg Config

	repo *pg.Repository
	srv  *grpc.Server
}

func NewDaemon(ctx context.Context, cfg Config) (*Daemon, error) {
	var ath auth.Authenticator
	if cfg.Auth.Disable {
		ath = auth.NewDisabledAuthenticator()
	} else {
		pk, err := auth.LoadPublicKey(cfg.Auth.PublicKeyFile)
		if err != nil {
			return nil, fmt.Errorf("load key: %w", err)
		}
		ath = auth.NewStaticAuthenticator(pk)
	}

	repo, err := pg.NewRepository(ctx, cfg.Database.DBS)
	if err != nil {
		return nil, fmt.Errorf("new repo: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logging.UnaryInterceptor(logging.L(ctx))),
	)
	standup.RegisterStandupServer(grpcServer, grpcserver.NewService(repo, ath))

	return &Daemon{
		cfg:  cfg,
		repo: repo,
		srv:  grpcServer,
	}, nil
}

func (d *Daemon) Start(ctx context.Context) error {
	logging.L(ctx).Sugar().Infof("Listening on %q...", d.cfg.Service.Addr)
	lis, err := net.Listen("tcp", d.cfg.Service.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return d.srv.Serve(lis)
}

func (d *Daemon) GracefulStop(ctx context.Context) error {
	logging.L(ctx).Sugar().Infof("Stopping the daemon...")
	d.srv.GracefulStop()
	return d.repo.Close(ctx)
}
