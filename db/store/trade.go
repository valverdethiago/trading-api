package store

import (
	"context"

	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// TradeStore interface with persistence operations for Trade
type TradeStore interface {
	CreateTrade(ctx context.Context, arg db.CreateTradeParams) (db.Trade, error)
	ListTradesByAccount(ctx context.Context, accountUUID uuid.UUID) ([]db.Trade, error)
	GetTradeByID(ctx context.Context, tradeUUID uuid.UUID) (db.Trade, error)
	UpdateTrade(ctx context.Context, arg db.UpdateTradeParams) (db.Trade, error)
	UpdateTradeStatus(ctx context.Context, arg db.UpdateTradeStatusParams) (db.Trade, error)
}

// DbTradeStore implementation of trade store executing operations against a real SQL database
type DbTradeStore struct {
	queries db.Querier
}

// NewDbTradeStore builds a new instance of db trade store
func NewDbTradeStore(queries db.Querier) TradeStore {
	return &DbTradeStore{
		queries: queries,
	}
}

// CreateTrade creates a new trade
func (tradeStore *DbTradeStore) CreateTrade(ctx context.Context, arg db.CreateTradeParams) (db.Trade, error) {
	return tradeStore.queries.CreateTrade(ctx, arg)
}

// ListTradesByAccount returns all trades of the given account
func (tradeStore *DbTradeStore) ListTradesByAccount(ctx context.Context, accountUUID uuid.UUID) ([]db.Trade, error) {
	return tradeStore.queries.ListTradesByAccount(ctx, accountUUID)
}

// GetTradeByID returns a trade with the given id
func (tradeStore *DbTradeStore) GetTradeByID(ctx context.Context, tradeUUID uuid.UUID) (db.Trade, error) {
	return tradeStore.queries.GetTradeById(ctx, tradeUUID)
}

// UpdateTrade updates a trade
func (tradeStore *DbTradeStore) UpdateTrade(ctx context.Context, arg db.UpdateTradeParams) (db.Trade, error) {
	return tradeStore.queries.UpdateTrade(ctx, arg)
}

// UpdateTradeStatus updates a trade status
func (tradeStore *DbTradeStore) UpdateTradeStatus(ctx context.Context, arg db.UpdateTradeStatusParams) (db.Trade, error) {
	return tradeStore.queries.UpdateTradeStatus(ctx, arg)
}
