package service

import (
	"context"
	"errors"
	"log"

	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// TradeService service to handle business rules for trades
type TradeService struct {
	queries        db.Querier
	accountService *AccountService
}

// NewTradeService creates a new TradeService instance
func NewTradeService(queries db.Querier, accountService *AccountService) *TradeService {
	return &TradeService{
		queries: queries,
	}
}

// CreateTrade Creates a new trade for the account
func (service *TradeService) CreateTrade(trade db.Trade, accountUUID string) (db.Trade, error) {
	var dbTrade db.Trade
	dbAccount, err := service.accountService.AssertAccountExists(accountUUID)
	if err != nil {
		return dbTrade, err
	}
	arg := db.CreateTradeParams{
		AccountUuid: dbAccount.AccountUuid,
		Symbol:      trade.Symbol,
		Quantity:    trade.Quantity,
		Side:        trade.Side,
		Price:       trade.Price,
	}
	return service.queries.CreateTrade(context.Background(), arg)
}

//ListTradesByAccount list all trades for a given account
func (service *TradeService) ListTradesByAccount(accountUUID string) ([]db.Trade, error) {
	var dbTrades []db.Trade
	dbAccount, err := service.accountService.AssertAccountExists(accountUUID)
	if err != nil {
		return dbTrades, err
	}
	return service.queries.ListTradesByAccount(context.Background(), dbAccount.AccountUuid)
}

//FindByIDAndAccountID finds a trade by its ID and account ID
func (service *TradeService) FindByIDAndAccountID(ID string, accountUUID string) (db.Trade, error) {
	return service.assertTradeExistsAndBelongToTheAccount(ID, accountUUID)
}

//CancelTradeByIDAndAccountID cancels a trade with the given id
func (service *TradeService) CancelTradeByIDAndAccountID(ID string, accountUUID string) (db.Trade, error) {
	dbTrade, err := service.assertTradeExistsAndBelongToTheAccount(ID, accountUUID)
	if err != nil {
		return dbTrade, err
	}
	if dbTrade.Status != db.TradeStatusSUBMITTED {
		return dbTrade, errors.New("It's not allowed to cancel a trade that are not on submitted status")
	}
	arg := db.UpdateTradeStatusParams{
		TradeUuid: dbTrade.TradeUuid,
		Status:    db.TradeStatusCANCELLED,
	}
	return service.queries.UpdateTradeStatus(context.Background(), arg)

}

// AssertTradeExists Returns the trade with the given ID
func (service *TradeService) AssertTradeExists(ID string) (db.Trade, error) {
	var dbTrade db.Trade
	uuid, err := parseUUID(ID)
	if err != nil {
		return dbTrade, err
	}
	log.Printf("Found trade with id %s", uuid)
	return dbTrade, nil
}

func (service *TradeService) assertTradeExistsAndBelongToTheAccount(ID string, accountUUID string) (db.Trade, error) {
	var dbTrade db.Trade
	dbAccount, err := service.accountService.AssertAccountExists(accountUUID)
	if err != nil {
		return dbTrade, err
	}
	dbTrade, err = service.AssertTradeExists(ID)
	if err != nil {
		return dbTrade, err
	}
	if dbAccount.AccountUuid != dbTrade.TradeUuid {
		return dbTrade, errors.New("The trade is not attached to the given account")
	}
	return dbTrade, err
}
