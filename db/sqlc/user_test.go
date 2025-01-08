package db

import (
	"context"
	"testing"
	"time"

	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt.Time, user2.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	ctx := context.Background()
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	updatedAccount, err := testQueries.UpdateUser(ctx, UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{String: newFullName, Valid: true},
	})

	require.NoError(t,err)
	require.NotEmpty(t, updatedAccount)
	require.NotEqual(t, oldUser.FullName, updatedAccount.FullName)
	require.Equal(t, newFullName, updatedAccount.FullName)

	require.Equal(t,oldUser.Email,updatedAccount.Email)
	require.Equal(t,oldUser.Username,updatedAccount.Username)
	require.Equal(t,oldUser.HashedPassword,updatedAccount.HashedPassword)


}

func TestUpdateUserOnlyEmail(t *testing.T) {
	ctx := context.Background()
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	updatedAccount, err := testQueries.UpdateUser(ctx, UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{String: newEmail, Valid: true},
	})

	require.NoError(t,err)
	require.NotEmpty(t, updatedAccount)
	require.NotEqual(t, oldUser.Email, updatedAccount.Email)
	require.Equal(t, newEmail, updatedAccount.Email)

	require.Equal(t,oldUser.FullName,updatedAccount.FullName)
	require.Equal(t,oldUser.Username,updatedAccount.Username)
	require.Equal(t,oldUser.HashedPassword,updatedAccount.HashedPassword)


}


func TestUpdateUserOnlyPassword(t *testing.T) {
	ctx := context.Background()
	oldUser := createRandomUser(t)
	newPassoword := util.RandomString(10)

	newHashedPassword,err := util.HashPassword(newPassoword)

	updatedAccount, err := testQueries.UpdateUser(ctx, UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{String: newHashedPassword, Valid: true},
	})

	require.NoError(t,err)
	require.NotEmpty(t, updatedAccount)
	require.NotEqual(t, oldUser.HashedPassword, updatedAccount.HashedPassword)
	require.Equal(t, newHashedPassword, updatedAccount.HashedPassword)

	require.Equal(t,oldUser.FullName,updatedAccount.FullName)
	require.Equal(t,oldUser.Username,updatedAccount.Username)
	require.Equal(t,oldUser.Email,updatedAccount.Email)


}
