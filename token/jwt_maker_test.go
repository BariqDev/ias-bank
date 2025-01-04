package token

import (
	"testing"
	"time"

	"github.com/BariqDev/ias-bank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt:= time.Now()
	expiredAt:= issuedAt.Add(duration)

	token, tokenPayload,err := maker.CreateToken(username, duration)
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

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, tokenPayload,err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, tokenPayload)

	payload, err := maker.VerifyToken(token)       
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTTokenAlg(t *testing.T) {
	username := util.RandomOwner()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	claims := &JWTPayload{
		payload,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        payload.ID.String(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)


	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)


	payload,err = maker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}