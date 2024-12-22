package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	dbSource = "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable"
)

var testQueries *Queries
var testDbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	testDbPool, err = pgxpool.New(ctx, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
		os.Exit(1)
	}
	testQueries = New(testDbPool)
	os.Exit(m.Run())
}
