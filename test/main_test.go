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
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/config"
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
}

func runTests(m *testing.M) error {
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

	conn, err := grpc.Dial(cfg.Standup.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
