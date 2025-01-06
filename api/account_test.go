package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/BariqDev/ias-bank/db/mock"
	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/token"
	"github.com/BariqDev/ias-bank/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, tokenMaker token.Maker, request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationBearerKey, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				checkResponseBody(t, recorder.Body, account)
			},
		},
		{
			name:      "notFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationBearerKey, user.Username, time.Minute)
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "internalServerError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationBearerKey, user.Username, time.Minute)
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "badRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(0)
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationBearerKey, user.Username, time.Minute)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build the stubs
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, server.tokenMaker, request)
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)

		})
	}

}

func randomAccount(owner string) db.Account {

	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

}

func checkResponseBody(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.Equal(t, account, gotAccount)
}
