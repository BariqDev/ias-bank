package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	ctx := context.Background()
	store := NewStore(testDbPool)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: pgtype.Int8{Int64: account1.ID, Valid: true},
				ToAccountID:   pgtype.Int8{Int64: account2.ID, Valid: true},
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		// check transfer
		transfer := result.Transfer
		require.Equal(t, int64(transfer.FromAccountID.Int64), account1.ID)
		require.Equal(t, int64(transfer.ToAccountID.Int64), account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(ctx, transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, int64(fromEntry.AccountID.Int64))
		require.Equal(t, -amount, int64(fromEntry.Amount.Int64))
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(ctx, fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, int64(toEntry.AccountID.Int64))
		require.Equal(t, amount, int64(toEntry.Amount.Int64))
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(ctx, toEntry.ID)
		require.NoError(t, err)
	}
}
