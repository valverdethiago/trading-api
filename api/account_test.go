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
	"github.com/stretchr/testify/require"
	mockdb "github.com/valverdethiago/trading-api/db/mock"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/util"
)

func TestCreateAccount(t *testing.T) {

	testCases := []struct {
		name          string
		buildRequest  func() CreateAccountRequest
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK Without Address",
			buildRequest: func() CreateAccountRequest {
				return CreateAccountRequest{
					Username: account.Username,
					Email:    account.Email,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
						Username: account.Username,
						Email:    account.Email,
					})).
					Times(1).
					Return(expectedAccount, nil)
				querier.EXPECT().
					GetAccountByUsername(gomock.Any(), gomock.Eq(account.Username)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, expectedAccount)
			},
		}, {
			name: "Internal Server Error",
			buildRequest: func() CreateAccountRequest {
				return CreateAccountRequest{
					Username: account.Username,
					Email:    account.Email,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountByUsername(gomock.Any(), account.Username).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				querier.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name: "Error Without Address And Username",
			buildRequest: func() CreateAccountRequest {
				return CreateAccountRequest{
					Email: account.Email,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "Error Without Address And Email",
			buildRequest: func() CreateAccountRequest {
				return CreateAccountRequest{
					Username: account.Username,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "OK With Address",
			buildRequest: func() CreateAccountRequest {
				addressReq := AddressRequest{
					Name:    expectedAddress.Name,
					Street:  expectedAddress.Street,
					City:    expectedAddress.City,
					State:   string(expectedAddress.State),
					Zipcode: expectedAddress.Zipcode,
				}
				return CreateAccountRequest{
					Username: account.Username,
					Email:    account.Email,
					Address:  &addressReq,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
						Username: account.Username,
						Email:    account.Email,
					})).
					Times(1).
					Return(expectedAccount, nil)
				querier.EXPECT().
					GetAccountByUsername(gomock.Any(), gomock.Eq(account.Username)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				querier.EXPECT().
					CreateAddress(gomock.Any(), gomock.Any()).
					Times(1).
					Return(expectedAddress, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, expectedAccount)
			},
		}, {
			name: "Error With Invalid Address",
			buildRequest: func() CreateAccountRequest {
				addressReq := AddressRequest{
					Street:  expectedAddress.Street,
					City:    expectedAddress.City,
					State:   string(expectedAddress.State),
					Zipcode: expectedAddress.Zipcode,
				}
				return CreateAccountRequest{
					Username: account.Username,
					Email:    account.Email,
					Address:  &addressReq,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
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
			request, err := http.NewRequest(http.MethodPost, "/accounts",
				sendObjectAsRequestBody(t, testCase.buildRequest()))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestListAccounts(t *testing.T) {
	testCases := []struct {
		name          string
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Empty List",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					ListAccounts(gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		}, {
			name: "OK",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					ListAccounts(gomock.Any()).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccountList(t, recorder.Body, accounts)
			},
		}, {
			name: "Internal Server Error",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					ListAccounts(gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			testCase.checkResponse(t, recorder)
		})

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

func requireBodyMatchAccountList(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var bodyAccounts []db.Account
	err = json.Unmarshal(data, &bodyAccounts)
	require.NoError(t, err)
	for _, account := range accounts {
		dbAccount := findAccountInList(bodyAccounts, account)
		require.Equal(t, dbAccount, account)
	}
}

func sendObjectAsRequestBody(t *testing.T, obj interface{}) *bytes.Buffer {
	b, err := json.Marshal(obj)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func findAccountInList(accounts []db.Account, account db.Account) db.Account {
	var result db.Account
	for _, element := range accounts {
		if element.AccountUuid == account.AccountUuid {
			result = element
		}
	}
	return result
}
