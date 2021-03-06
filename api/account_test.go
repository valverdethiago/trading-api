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
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mockdb "github.com/valverdethiago/trading-api/db/mock"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/util"
)

func createRandomAccount() db.Account {
	return db.Account{
		AccountUuid: uuid.New(),
		Username:    util.RandomUsername(),
		Email:       util.RandomEmail(),
	}
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount()

	testCases := []struct {
		name          string
		id            string
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   account.AccountUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		}, {
			name: "NOT FOUND",
			id:   account.AccountUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name: "DATABASE CONNECTION LOST",
			id:   account.AccountUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name: "INVALID ID",
			id:   util.RandomNumericString(4),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mockdb.NewMockQuerier(ctrl)
			// build stubs
			testCase.buildStubs(querier)

			// start http server and send the request
			server := NewServer(querier)
			recorder := httptest.NewRecorder()
			urlToTest := fmt.Sprintf("/accounts/%s", testCase.id)
			request, err := http.NewRequest(http.MethodGet, urlToTest, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			testCase.checkResponse(t, recorder)
		})

	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var bodyAccount db.Account
	err = json.Unmarshal(data, &bodyAccount)
	require.NoError(t, err)
	require.Equal(t, account, bodyAccount)
}
