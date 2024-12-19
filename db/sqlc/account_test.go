package db

import (
	"context"
	"testing"
	"time"

	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	ctx := context.Background()
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(ctx, args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Currency, args.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAcc(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccount(ctx, account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)

}
func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(ctx, args)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(ctx, account1.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(ctx, account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account)
}


func TestListAccount(t *testing.T) {
	
}