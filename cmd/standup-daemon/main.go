package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/xopoww/standup/internal/grpcserver"
	"github.com/xopoww/standup/internal/repository/pg"
	"github.com/xopoww/standup/internal/repository/pg/migrations"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Uint("port", 65000, "tcp port")
	dbs := flag.String("dbs", "", "DBS")
	flag.Parse()

	if *dbs == "" {
		log.Fatal("dbs must be specified")
	}

	ctx := context.Background()

	mig, err := migrations.NewMigration(ctx, *dbs)
	if err != nil {
		log.Fatal(err)
	}
	err = mig.Up()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := pg.NewRepository(ctx, *dbs)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	standup.RegisterStandupServer(grpcServer, grpcserver.NewService(repo))
	log.Fatal(grpcServer.Serve(lis))
}
