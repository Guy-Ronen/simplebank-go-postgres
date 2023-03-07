package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/guy-ronen/simplebank/db/mock"
	db "github.com/guy-ronen/simplebank/db/sqlc"
	"github.com/guy-ronen/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := RandomAccount()

	// create a test case struct to hold the test case data and expected results for each test case
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// test case 1 - get account successfully
		{
			name:      "OK",
			accountID: account.ID,

			// build the stubs for the mock store
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},

		// test case 2 - get account not found
		{
			name:      "NotFound",
			accountID: account.ID,

			// build the stubs for the mock store
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		// test case 3 -  Internal Server Error
		{
			name:      "InternalError",
			accountID: account.ID,

			// build the stubs for the mock store
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// test case 4 -  Invalid ID
		{
			name:      "InvalidID",
			accountID: 0,

			// build the stubs for the mock store
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// loop through the test cases
	for i := range testCases {

		// get the test case
		tc := testCases[i]

		// run the test case
		t.Run(tc.name, func(t *testing.T) {

			// create a new mock controller
			ctrl := gomock.NewController(t)

			// defer the call to Finish() to clean up the controller
			defer ctrl.Finish()

			// create a mock store
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// create a new server
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// create a new request
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// send request to the server
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// TODO - implemenet createAccountAPI test and DeleteAccountAPI test

func RandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {

	// read the response body
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// unmarshal the response body
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
