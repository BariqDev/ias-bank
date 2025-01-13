package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// query provide all functions to execute db and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, args CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	*Queries
	pool *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(pool),
		pool:    pool,
	}
}


// execTx executes function with database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil {

		if rbErr := tx.Rollback(ctx); rbErr != nil {

			return fmt.Errorf("tx error %v rollback error %v", err, rbErr)
		}
		return err

	}

	return tx.Commit(ctx)
}

