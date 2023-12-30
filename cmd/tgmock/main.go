package main

import (
	"flag"

	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/pkg/config"
	"github.com/xopoww/standup/pkg/tgmock"
)

func main() {
	var args struct {
		cfgPath string
	}

	logger := logging.NewLogger()
	defer func() {
		_ = logger.Sync()
	}()

	flag.StringVar(&args.cfgPath, "config", "", "path to yaml config file")
	flag.Parse()
	if args.cfgPath == "" {
		logger.Sugar().Fatal("`config` is required.")
	}

	var cfg tgmock.Config
	err := config.LoadFile(args.cfgPath, &cfg)
	if err != nil {
		logger.Sugar().Fatalf("Load config: %s.", err)
	}

	tm := tgmock.New(cfg, logger)
	logger.Sugar().Errorf("TGMock stopped: %s.", tm.Start())
}
