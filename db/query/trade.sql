-- name: GetTradeById :one
SELECT * 
  FROM trade
 WHERE trade_uuid = $1;

-- name: ListTradesByAccount :many
  SELECT * 
    FROM trade
   WHERE account_uuid = $1
ORDER BY created_date;

-- name: CreateTrade :one
INSERT INTO trade (account_uuid, symbol, quantity, side          , price) 
     VALUES       ($1          , $2    , $3      , $4::trade_side, $5   )
RETURNING *; 

-- name: UpdateTrade :one
UPDATE trade 
   SET symbol = $1, 
       quantity = $2,
       side = $3::trade_side,
       price = $4, 
       status = $5::trade_status,
       updated_date = now()
 WHERE trade_uuid = $6
 RETURNING *;