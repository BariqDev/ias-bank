package token

import (
	"testing"
	"time"

	"github.com/BariqDev/ias-bank/util"
	"github.com/stretchr/testify/require"
)


func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt:= time.Now()
	expiredAt:= issuedAt.Add(duration)

	token,tokenPayload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, tokenPayload)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t,username,payload.Username)
	require.NotZero(t, payload.ID)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, tokenPayload,err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, tokenPayload)

	payload, err := maker.VerifyToken(token)       
	require.Error(t, err)
	require.Nil(t, payload)

}

