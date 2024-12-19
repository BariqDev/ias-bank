package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	testQueries = New(conn)
	os.Exit(m.Run())
}
