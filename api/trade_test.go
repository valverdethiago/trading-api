package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func TestCreateTrade(t *testing.T) {
	testCases := []struct {
		name          string
		accountID     string
		buildRequest  func() tradeRequest
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.AccountUuid.String(),
			buildRequest: func() tradeRequest {
				return tradeRequest{
					Symbol:   trade.Symbol,
					Quantity: trade.Quantity,
					Side:     trade.Side,
					Price:    trade.Price,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					CreateTrade(gomock.Any(), gomock.Any()).
					Return(trade, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTrade(t, recorder.Body, trade)
			},
		}, {
			name:      "Missing account ID",
			accountID: "",
			buildRequest: func() tradeRequest {
				return tradeRequest{
					Symbol:   trade.Symbol,
					Quantity: trade.Quantity,
					Side:     trade.Side,
					Price:    trade.Price,
				}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Invalid Trade Request",
			accountID: account.AccountUuid.String(),
			buildRequest: func() tradeRequest {
				return tradeRequest{}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Invalid AccountID",
			accountID: util.RandomAlphaNumericString(4),
			buildRequest: func() tradeRequest {
				return tradeRequest{
					Symbol:   trade.Symbol,
					Quantity: trade.Quantity,
					Side:     trade.Side,
					Price:    trade.Price}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Unexistent Account",
			accountID: account.AccountUuid.String(),
			buildRequest: func() tradeRequest {
				return tradeRequest{
					Symbol:   trade.Symbol,
					Quantity: trade.Quantity,
					Side:     trade.Side,
					Price:    trade.Price}
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
			name:      "Internal Server Error",
			accountID: account.AccountUuid.String(),
			buildRequest: func() tradeRequest {
				return tradeRequest{
					Symbol:   trade.Symbol,
					Quantity: trade.Quantity,
					Side:     trade.Side,
					Price:    trade.Price}
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		querier := mockdb.NewMockQuerier(ctrl)
		//build stubs
		testCase.buildStubs(querier)

		// start server and send the request
		server := NewServer(querier)
		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/accounts/%s/trades", testCase.accountID)
		request, err := http.NewRequest(http.MethodPost, url, sendObjectAsRequestBody(t, testCase.buildRequest()))
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		//check response
		testCase.checkResponse(t, recorder)
	}
}

func TestListTrade(t *testing.T) {
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
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					ListTradesByAccount(gomock.Any(), account.AccountUuid).
					Return(trades, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTradeList(t, recorder.Body, trades)
			},
		}, {
			name:       "Missing account ID",
			accountID:  "",
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "Invalid AccountID",
			accountID:  util.RandomAlphaNumericString(4),
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Unexistent Account",
			accountID: account.AccountUuid.String(),
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
			name:      "Internal Server Error",
			accountID: account.AccountUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		querier := mockdb.NewMockQuerier(ctrl)
		//build stubs
		testCase.buildStubs(querier)

		// start server and send the request
		server := NewServer(querier)
		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/accounts/%s/trades", testCase.accountID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		//check response
		testCase.checkResponse(t, recorder)
	}
}

func requireBodyMatchTrade(t *testing.T, body *bytes.Buffer, trade db.Trade) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var bodyTrade db.Trade
	err = json.Unmarshal(data, &bodyTrade)
	require.NoError(t, err)
	require.Equal(t, trade, bodyTrade)
}

func requireBodyMatchTradeList(t *testing.T, body *bytes.Buffer, trades []db.Trade) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var bodyTrades []db.Trade
	err = json.Unmarshal(data, &bodyTrades)
	require.NoError(t, err)
	for _, trade := range trades {
		dbTrade := findTradeInList(bodyTrades, trade)
		require.Equal(t, dbTrade, trade)
	}
}

func findTradeInList(trades []db.Trade, trade db.Trade) db.Trade {
	var result db.Trade
	for _, element := range trades {
		if element.TradeUuid == trade.TradeUuid {
			result = element
		}
	}
	return result
}

func TestCancelTrade(t *testing.T) {
	testCases := []struct {
		name          string
		accountID     string
		tradeID       string
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.AccountUuid.String(),
			tradeID:   trade.TradeUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetTradeById(gomock.Any(), gomock.Eq(trade.TradeUuid)).
					Times(1).
					Return(expectedSubmittedTrade, nil)
				querier.EXPECT().
					UpdateTradeStatus(gomock.Any(), gomock.Eq(db.UpdateTradeStatusParams{
						TradeUuid: trade.TradeUuid,
						Status:    db.TradeStatusCANCELLED,
					})).
					Return(expectedCanceledTrade, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				requireBodyMatchTrade(t, recorder.Body, expectedCanceledTrade)
			},
		}, {
			name:      "Trade doesn't belong to the account",
			accountID: account.AccountUuid.String(),
			tradeID:   trade.TradeUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetTradeById(gomock.Any(), gomock.Eq(trade.TradeUuid)).
					Times(1).
					DoAndReturn(func(ctx context.Context, tradeUuid uuid.UUID) (db.Trade, error) {
						return db.Trade{
							TradeUuid:   trade.TradeUuid,
							AccountUuid: uuid.New(),
							Symbol:      trade.Symbol,
							Quantity:    trade.Quantity,
							Side:        trade.Side,
							Price:       trade.Price,
							Status:      trade.Status,
						}, nil
					})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:      "OK",
			accountID: account.AccountUuid.String(),
			tradeID:   trade.TradeUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetTradeById(gomock.Any(), gomock.Eq(trade.TradeUuid)).
					Times(1).
					Return(expectedSubmittedTrade, nil)
				querier.EXPECT().
					UpdateTradeStatus(gomock.Any(), gomock.Eq(db.UpdateTradeStatusParams{
						TradeUuid: trade.TradeUuid,
						Status:    db.TradeStatusCANCELLED,
					})).
					Return(expectedCanceledTrade, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				requireBodyMatchTrade(t, recorder.Body, expectedCanceledTrade)
			},
		}, {
			name:      "Already Cancelled",
			accountID: account.AccountUuid.String(),
			tradeID:   trade.TradeUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.AccountUuid)).
					Times(1).
					Return(account, nil)
				querier.EXPECT().
					GetTradeById(gomock.Any(), gomock.Eq(trade.TradeUuid)).
					Times(1).
					Return(expectedCanceledTrade, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		}, {
			name:       "Missing account ID",
			accountID:  "",
			tradeID:    trade.TradeUuid.String(),
			buildStubs: func(querier *mockdb.MockQuerier) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		querier := mockdb.NewMockQuerier(ctrl)
		//build stubs
		testCase.buildStubs(querier)

		// start server and send the request
		server := NewServer(querier)
		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/accounts/%s/trades/%s", testCase.accountID, testCase.tradeID)
		log.Print(url)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		//check response
		testCase.checkResponse(t, recorder)
	}
}
