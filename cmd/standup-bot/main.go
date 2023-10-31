package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/xopoww/standup/internal/auth"
	"github.com/xopoww/standup/internal/bot"
	"github.com/xopoww/standup/internal/bot/repository/pg"
	"github.com/xopoww/standup/internal/bot/service"
	"github.com/xopoww/standup/internal/bot/tg"
	"github.com/xopoww/standup/internal/logging"
	"github.com/xopoww/standup/pkg/api/standup"
	"github.com/xopoww/standup/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var args struct {
		cfgPath string
		devel   bool
	}

	logger := logging.NewLogger()
	defer func() {
		_ = logger.Sync()
	}()

	flag.StringVar(&args.cfgPath, "config", "", "path to yaml config file")
	flag.BoolVar(&args.devel, "devel", false, "enable development mode")
	flag.Parse()
	if args.cfgPath == "" {
		logger.Sugar().Fatal("`config` is required.")
	}

	var cfg bot.Config
	err := config.LoadFile(args.cfgPath, &cfg)
	if err != nil {
		logger.Sugar().Fatalf("Load config: %s.", err)
	}

	ctx := logging.WithLogger(context.Background(), logger)

	repo, err := pg.NewRepository(ctx, cfg.Database.DBS)
	if err != nil {
		logger.Sugar().Fatalf("Init repository: %s.", err)
	}

	tgBot, err := tg.NewBot(ctx, *cfg.Bot, args.devel)
	if err != nil {
		logger.Sugar().Fatalf("Init telegram bot: %s.", err)
	}

	conn, err := grpc.DialContext(ctx, cfg.Standup.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Sugar().Fatalf("Dial GRPC: %s.", err)
	}
	client := standup.NewStandupClient(conn)

	pk, err := auth.LoadPrivateKey(cfg.PrivateKeyFile)
	if err != nil {
		logger.Sugar().Fatalf("Load private key: %s.", err)
	}
	srv, err := service.NewService(logger, *cfg.Service, service.Deps{
		Bot:    tgBot,
		Repo:   repo,
		Client: client,
		Issuer: auth.NewStaticIssuer(pk),
	})
	if err != nil {
		logger.Sugar().Fatalf("New service: %s.")
	}

	srv.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sig := <-ch
	logger.Sugar().Debugf("Received signal %s, stopping the service...", sig)
	srv.Stop()
	logger.Sugar().Infof("Stopped.")
}
