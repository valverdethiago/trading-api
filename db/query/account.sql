-- name: GetAccountByUsername :one
SELECT * 
  FROM account
 WHERE username = $1;

-- name: GetAccountById :one
SELECT * 
  FROM account
 WHERE account_uuid = $1;

-- name: ListAccounts :many
  SELECT * 
    FROM account
ORDER BY created_date;

-- name: CreateAccount :one
INSERT INTO account (username, email) 
VALUES ($1, $2)
RETURNING *; 

-- name: UpdateAccount :one
UPDATE account 
   SET username = $1, 
       email = $2,
       updated_date = now()
 WHERE account_uuid = $3
 RETURNING *;