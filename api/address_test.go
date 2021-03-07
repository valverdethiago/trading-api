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

func TestCreateAddressForAccount(t *testing.T) {
	testCases := []struct {
		name          string
		accountID     string
		buildRequest  func() AddressRequest
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(expectedAccount, nil)
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Address{}, sql.ErrNoRows)
				querier.EXPECT().
					CreateAddress(gomock.Any(), gomock.Any()).
					Times(1).
					Return(expectedAddress, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAddress(t, recorder.Body, expectedAddress)
			},
		}, {
			name:      "InvalidID",
			accountID: util.RandomAlphaNumericString(4),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "InvalidAddress",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Account Doesn't exist",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
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
			name:      "Account Already Have Address",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(address, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
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
			url := fmt.Sprintf("/accounts/%s/address", testCase.accountID)
			request, err := http.NewRequest(http.MethodPut, url,
				sendObjectAsRequestBody(t, testCase.buildRequest()))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestUpdateAddressForAccount(t *testing.T) {
	testCases := []struct {
		name          string
		accountID     string
		buildRequest  func() AddressRequest
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(expectedAccount, nil)
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(address, nil)
				querier.EXPECT().
					UpdateAddress(gomock.Any(), gomock.Any()).
					Times(1).
					Return(expectedAddress, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAddress(t, recorder.Body, expectedAddress)
			},
		}, {
			name:      "InvalidID",
			accountID: util.RandomAlphaNumericString(4),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "InvalidAddress",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Account Doesn't exist",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
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
			name:      "Account Doesn't Have Address",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Address{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:      "Internal Server Error",
			accountID: account.AccountUuid.String(),
			buildRequest: func() AddressRequest {
				return AddressRequest{
					Name:    address.Name,
					Street:  address.Street,
					City:    address.City,
					State:   string(address.State),
					Zipcode: address.Zipcode,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(address, nil)
				querier.EXPECT().
					UpdateAddress(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Address{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/accounts/%s/address", testCase.accountID)
			request, err := http.NewRequest(http.MethodPost, url,
				sendObjectAsRequestBody(t, testCase.buildRequest()))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			testCase.checkResponse(t, recorder)
		})
	}
}
func TestGetAddressForAccount(t *testing.T) {
	testCases := []struct {
		name          string
		accountID     string
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.AccountUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(expectedAddress, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAddress(t, recorder.Body, expectedAddress)
			},
		}, {
			name:       "InvalidID",
			accountID:  util.RandomAlphaNumericString(4),
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Account Doesn't Have Address",
			accountID: account.AccountUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAddressByAccount(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Address{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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
			url := fmt.Sprintf("/accounts/%s/address", testCase.accountID)
			request, err := http.NewRequest(http.MethodGet, url,
				sendObjectAsRequestBody(t, nil))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			testCase.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchAddress(t *testing.T, body *bytes.Buffer, address db.Address) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var bodyAddress db.Address
	err = json.Unmarshal(data, &bodyAddress)
	require.NoError(t, err)
	require.Equal(t, address, bodyAddress)
}
