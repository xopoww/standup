package test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/xopoww/standup/internal/logging"
	"github.com/xopoww/standup/internal/testutil"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Standup struct {
		Addr string `yaml:"addr" validate:"required,hostname_port"`
	} `yaml:"standup"`
	Database struct {
		DBS string `yaml:"dbs" validate:"required"`
	}
}

var args struct {
	cfgPath string
}

var deps struct {
	client standup.StandupClient
	db     *pgx.Conn
	logger *zap.Logger
}

func runTests(m *testing.M) error {
	deps.logger = logging.NewLogger()
	defer deps.logger.Sync()

	flag.StringVar(&args.cfgPath, "config", "", "path to yaml config file")
	flag.Parse()
	if args.cfgPath == "" {
		return errors.New("`config` must be specified")
	}

	var cfg Config
	err := config.LoadFile(args.cfgPath, &cfg)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	conn, err := grpc.Dial(cfg.Standup.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(logging.UnaryClientInterceptor(deps.logger)),
	)
	if err != nil {
		return fmt.Errorf("grpc dial: %w", err)
	}
	deps.client = standup.NewStandupClient(conn)

	db, err := pgx.Connect(context.TODO(), cfg.Database.DBS)
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}
	defer db.Close(context.TODO())
	deps.db = db

	if rc := m.Run(); rc != 0 {
		log.Printf("Tests finished with code %d.", rc)
		os.Exit(rc)
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := runTests(m); err != nil {
		log.Fatalf("Fatal error: %s.", err)
	}
}

func RunTest(t *testing.T, name string, f func(context.Context, *testing.T)) {
	ctx, cancel := testutil.NewContext(context.Background())
	defer cancel()
	ctx = logging.WithLogger(ctx, deps.logger)
	t.Run(name, func(tt *testing.T) {
		logging.L(ctx).Sugar().Infof("Running %s with ID %q...", t.Name(), testutil.TestID(ctx))
		f(ctx, tt)
	})
}
