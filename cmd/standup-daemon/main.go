package main

import (
	"context"
	"flag"

	"github.com/xopoww/standup/internal/daemon"
	"github.com/xopoww/standup/internal/logging"
	"github.com/xopoww/standup/pkg/config"
)

var args struct {
	cfgPath string
}

func main() {
	logger := logging.NewLogger()
	defer logger.Sync()

	flag.StringVar(&args.cfgPath, "config", "", "path to yaml config file")
	flag.Parse()
	if args.cfgPath == "" {
		logger.Sugar().Fatal("`config` is required.")
	}

	var cfg daemon.Config
	err := config.LoadFile(args.cfgPath, &cfg)
	if err != nil {
		logger.Sugar().Fatalf("Load config: %s.", err)
	}

	ctx := logging.WithLogger(context.Background(), logger)
	d, err := daemon.NewDaemon(ctx, cfg)
	if err != nil {
		logger.Sugar().Fatalf("Init daemon: %s.", err)
	}

	if err := d.Start(ctx); err != nil {
		logger.Sugar().Fatalf("Daemon terminated with error: %s.", err)
	} else {
		logger.Sugar().Info("Daemon stopped.")
	}
}
