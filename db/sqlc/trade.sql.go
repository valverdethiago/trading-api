// Code generated by sqlc. DO NOT EDIT.
// source: trade.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createTrade = `-- name: CreateTrade :one
INSERT INTO trade (account_uuid, symbol, quantity, side          , price) 
     VALUES       ($1          , $2    , $3      , $4::trade_side, $5   )
RETURNING trade_uuid, account_uuid, symbol, quantity, side, price, status, created_date, updated_date, created_by, updated_by
`

type CreateTradeParams struct {
	AccountUuid uuid.UUID `json:"account_uuid"`
	Symbol      string    `json:"symbol"`
	Quantity    int64     `json:"quantity"`
	Side        TradeSide `json:"side"`
	Price       float64   `json:"price"`
}

func (q *Queries) CreateTrade(ctx context.Context, arg CreateTradeParams) (Trade, error) {
	row := q.db.QueryRowContext(ctx, createTrade,
		arg.AccountUuid,
		arg.Symbol,
		arg.Quantity,
		arg.Side,
		arg.Price,
	)
	var i Trade
	err := row.Scan(
		&i.TradeUuid,
		&i.AccountUuid,
		&i.Symbol,
		&i.Quantity,
		&i.Side,
		&i.Price,
		&i.Status,
		&i.CreatedDate,
		&i.UpdatedDate,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}

const getTradeById = `-- name: GetTradeById :one
SELECT trade_uuid, account_uuid, symbol, quantity, side, price, status, created_date, updated_date, created_by, updated_by 
  FROM trade
 WHERE trade_uuid = $1
`

func (q *Queries) GetTradeById(ctx context.Context, tradeUuid uuid.UUID) (Trade, error) {
	row := q.db.QueryRowContext(ctx, getTradeById, tradeUuid)
	var i Trade
	err := row.Scan(
		&i.TradeUuid,
		&i.AccountUuid,
		&i.Symbol,
		&i.Quantity,
		&i.Side,
		&i.Price,
		&i.Status,
		&i.CreatedDate,
		&i.UpdatedDate,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}

const listTradesByAccount = `-- name: ListTradesByAccount :many
  SELECT trade_uuid, account_uuid, symbol, quantity, side, price, status, created_date, updated_date, created_by, updated_by 
    FROM trade
   WHERE account_uuid = $1
ORDER BY created_date
`

func (q *Queries) ListTradesByAccount(ctx context.Context, accountUuid uuid.UUID) ([]Trade, error) {
	rows, err := q.db.QueryContext(ctx, listTradesByAccount, accountUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Trade
	for rows.Next() {
		var i Trade
		if err := rows.Scan(
			&i.TradeUuid,
			&i.AccountUuid,
			&i.Symbol,
			&i.Quantity,
			&i.Side,
			&i.Price,
			&i.Status,
			&i.CreatedDate,
			&i.UpdatedDate,
			&i.CreatedBy,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTrade = `-- name: UpdateTrade :one
UPDATE trade 
   SET symbol = $1, 
       quantity = $2,
       side = $3::trade_side,
       price = $4, 
       status = $5::trade_status,
       updated_date = now()
 WHERE trade_uuid = $6
 RETURNING trade_uuid, account_uuid, symbol, quantity, side, price, status, created_date, updated_date, created_by, updated_by
`

type UpdateTradeParams struct {
	Symbol    string      `json:"symbol"`
	Quantity  int64       `json:"quantity"`
	Side      TradeSide   `json:"side"`
	Price     float64     `json:"price"`
	Status    TradeStatus `json:"status"`
	TradeUuid uuid.UUID   `json:"trade_uuid"`
}

func (q *Queries) UpdateTrade(ctx context.Context, arg UpdateTradeParams) (Trade, error) {
	row := q.db.QueryRowContext(ctx, updateTrade,
		arg.Symbol,
		arg.Quantity,
		arg.Side,
		arg.Price,
		arg.Status,
		arg.TradeUuid,
	)
	var i Trade
	err := row.Scan(
		&i.TradeUuid,
		&i.AccountUuid,
		&i.Symbol,
		&i.Quantity,
		&i.Side,
		&i.Price,
		&i.Status,
		&i.CreatedDate,
		&i.UpdatedDate,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}
