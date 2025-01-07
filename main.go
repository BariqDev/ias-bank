package main

import (
	"context"
	"log"
	"net"
	"os"
	"github.com/BariqDev/ias-bank/api"
	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/gapi"
	"github.com/BariqDev/ias-bank/pb"
	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Load configuration file
	config, err := util.LoadConfig(".")
	ctx := context.Background()
	if err != nil {
		log.Fatal("Cannot load config:", err)
		os.Exit(1)
	}

	testDbPool, err := pgxpool.New(ctx, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
		os.Exit(1)
	}
	store := db.NewStore(testDbPool)

	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIASBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}


}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create server:", err)
		os.Exit(1)
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
		os.Exit(1)
	}

}
