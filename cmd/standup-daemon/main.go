package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/xopoww/standup/internal/grpcserver"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Uint("port", 65000, "tcp port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	standup.RegisterStandupServer(grpcServer, grpcserver.NewService())
	log.Fatal(grpcServer.Serve(lis))
}
