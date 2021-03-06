package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valverdethiago/trading-api/util"
)

func createRandomTrade(t *testing.T, account Account) Trade {
	arg := CreateTradeParams{
		AccountUuid: account.AccountUuid,
		Symbol:      util.RandomString(3),
		Quantity:    util.RandomInt(1, 1000),
		Side:        TradeSideBUY,
		Price:       util.RandomFloat(1, 1000),
	}

	trade, err := testQueries.CreateTrade(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trade)
	require.Equal(t, arg.AccountUuid, trade.AccountUuid)
	require.Equal(t, arg.Symbol, trade.Symbol)
	require.Equal(t, arg.Quantity, trade.Quantity)
	require.Equal(t, arg.Side, trade.Side)
	require.Equal(t, arg.Price, trade.Price)
	require.Equal(t, trade.Status, TradeStatusSUBMITTED)
	return trade
}

func TestCreateTrade(t *testing.T) {
	account := createRandomAccount(t)
	createRandomTrade(t, account)
}

func TestListTradesFromAccount(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTrade(t, account)
	}

	trades, err := testQueries.ListTradesByAccount(context.Background(), account.AccountUuid)
	require.NoError(t, err)
	require.NotEmpty(t, trades)
	require.Equal(t, len(trades), 5)
	for _, trade := range trades {
		require.NotEmpty(t, trade)
		require.NotEmpty(t, trade.AccountUuid)
		require.NotEmpty(t, trade.CreatedDate)
		require.NotEmpty(t, trade.Price)
		require.NotEmpty(t, trade.Quantity)
		require.NotEmpty(t, trade.Side)
		require.NotEmpty(t, trade.Status)
		require.Equal(t, trade.Status, TradeStatusSUBMITTED)
		require.NotEmpty(t, trade.Symbol)
		require.NotEmpty(t, trade.TradeUuid)
	}
}
func TestGetTradeById(t *testing.T) {
	account := createRandomAccount(t)
	trade := createRandomTrade(t, account)

	dbTrade, err := testQueries.GetTradeById(context.Background(), trade.TradeUuid)
	require.NoError(t, err)
	require.NotEmpty(t, dbTrade)
	require.Equal(t, trade.AccountUuid, dbTrade.AccountUuid)
	require.Equal(t, trade.CreatedDate, dbTrade.CreatedDate)
	require.Equal(t, trade.Price, dbTrade.Price)
	require.Equal(t, trade.Quantity, dbTrade.Quantity)
	require.Equal(t, trade.Side, dbTrade.Side)
	require.Equal(t, trade.Status, dbTrade.Status)
	require.Equal(t, trade.Status, dbTrade.Status)
	require.Equal(t, trade.Symbol, dbTrade.Symbol)
	require.Equal(t, trade.TradeUuid, dbTrade.TradeUuid)
}

func TestUpdateTrade(t *testing.T) {
	account := createRandomAccount(t)
	trade := createRandomTrade(t, account)

	arg := UpdateTradeParams{
		Symbol:    util.RandomAlphaNumericString(3),
		Side:      TradeSideSELL,
		Quantity:  trade.Quantity + 1,
		Price:     trade.Price + 1,
		Status:    TradeStatusFAILED,
		TradeUuid: trade.TradeUuid,
	}

	dbTrade, err := testQueries.UpdateTrade(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dbTrade)
	require.Equal(t, trade.AccountUuid, dbTrade.AccountUuid)
	require.Equal(t, trade.CreatedDate, dbTrade.CreatedDate)
	require.NotEqual(t, trade.UpdatedDate, dbTrade.UpdatedDate)
	require.NotEqual(t, trade.Symbol, dbTrade.Symbol)
	require.NotEqual(t, trade.Quantity, dbTrade.Quantity)
	require.NotEqual(t, trade.Side, dbTrade.Side)
	require.Equal(t, trade.TradeUuid, dbTrade.TradeUuid)
}
func TestUpdateTradeStatus(t *testing.T) {
	account := createRandomAccount(t)
	trade := createRandomTrade(t, account)

	arg := UpdateTradeStatusParams{
		Status:    TradeStatusCANCELLED,
		TradeUuid: trade.TradeUuid,
	}

	dbTrade, err := testQueries.UpdateTradeStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dbTrade)
	require.Equal(t, trade.AccountUuid, dbTrade.AccountUuid)
	require.Equal(t, trade.CreatedDate, dbTrade.CreatedDate)
	require.NotEqual(t, trade.UpdatedDate, dbTrade.UpdatedDate)
	require.Equal(t, trade.Symbol, dbTrade.Symbol)
	require.Equal(t, trade.Quantity, dbTrade.Quantity)
	require.Equal(t, trade.Side, dbTrade.Side)
	require.Equal(t, trade.TradeUuid, dbTrade.TradeUuid)
	require.Equal(t, dbTrade.Status, TradeStatusCANCELLED)
}
