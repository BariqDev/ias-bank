package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BariqDev/ias-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorizationHeader(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, tokenPayload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenPayload)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
func TestAuthMiddleware(t *testing.T) {

	testcases := []struct {
		name          string
		setupAuth     func(t *testing.T, tokenMaker token.Maker, request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationBearerKey, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "AuthorizationHeaderMissing",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidFormat",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, "", "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "unSupportedType",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, "unSupported", "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "expiredToken",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationBearerKey, "user", -time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker),
				authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			tc.setupAuth(t, server.tokenMaker, request)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
