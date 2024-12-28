package main

import (
	"context"
	"log"
	"os"

	"github.com/BariqDev/ias-bank/api"
	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)



func main() {

	// Load configuration file
	config,err := util.LoadConfig(".")
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

	server,err := api.NewServer(config,store)

	if err != nil {
		log.Fatal("Cannot create server:", err)
		os.Exit(1)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
		os.Exit(1)
	}

}
