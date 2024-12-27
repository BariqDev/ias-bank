package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)



var testQueries *Queries
var testDbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	config,err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load config:", err)
		os.Exit(1)
	}
	

	testDbPool, err = pgxpool.New(ctx, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
		os.Exit(1)
	}
	testQueries = New(testDbPool)
	os.Exit(m.Run())
}
