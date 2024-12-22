package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// query provide all functions to execute db and transactions
type Store struct {
	*Queries
	pool *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(pool),
		pool:    pool,
	}
}

type TransferTxParams struct {
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	Amount        int64       `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// execTx executes function with database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// TransferTx performs a money transfer from one account to the other
// It creates the transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		// create entry records
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    pgtype.Int8{Int64: -args.Amount, Valid: true},
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    pgtype.Int8{Int64: args.Amount, Valid: true},
		})
		if err != nil {
			return err
		}

		// TODO: update accounts' balance

		return nil
	})

	return result, err
}
