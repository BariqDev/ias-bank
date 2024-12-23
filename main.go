package main

import (
	"context"
	"log"
	"os"

	"github.com/BariqDev/ias-bank/api"
	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {

	ctx := context.Background()
	testDbPool, err := pgxpool.New(ctx, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
		os.Exit(1)
	}
	store := db.NewStore(testDbPool)

	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("Cannot start server:", err)
		os.Exit(1)
	}

}
