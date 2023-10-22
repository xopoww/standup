package main

import (
	"context"
	"flag"
	"log"

	"github.com/xopoww/standup/internal/daemon"
	"github.com/xopoww/standup/pkg/config"
)

var args struct {
	cfgPath string
}

func main() {
	flag.StringVar(&args.cfgPath, "config", "", "path to yaml config file")
	flag.Parse()
	if args.cfgPath == "" {
		log.Fatal("`config` must be specified")
	}

	var cfg daemon.Config
	err := config.LoadFile(args.cfgPath, &cfg)
	if err != nil {
		log.Fatalf("Load config: %s.", err)
	}

	ctx := context.Background()
	d, err := daemon.NewDaemon(ctx, cfg)
	if err != nil {
		log.Fatalf("Init daemon: %s.", err)
	}

	if err := d.Start(); err != nil {
		log.Fatalf("Daemon terminated with error: %s.", err)
	} else {
		log.Printf("Daemon stopped.")
	}
}
