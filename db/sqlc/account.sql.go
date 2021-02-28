// Code generated by sqlc. DO NOT EDIT.
// source: account.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO account (username, email) 
VALUES ($1, $2)
RETURNING account_uuid, username, email, created_date, updated_date, created_by, updated_by
`

type CreateAccountParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Username, arg.Email)
	var i Account
	err := row.Scan(
		&i.AccountUuid,
		&i.Username,
		&i.Email,
		&i.CreatedDate,
		&i.UpdatedDate,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}

const getAccountById = `-- name: GetAccountById :one
SELECT account_uuid, username, email, created_date, updated_date, created_by, updated_by 
  FROM account
 WHERE account_uuid = $1
`

func (q *Queries) GetAccountById(ctx context.Context, accountUuid uuid.UUID) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountById, accountUuid)
	var i Account
	err := row.Scan(
		&i.AccountUuid,
		&i.Username,
		&i.Email,
		&i.CreatedDate,
		&i.UpdatedDate,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
  SELECT account_uuid, username, email, created_date, updated_date, created_by, updated_by 
    FROM account
ORDER BY created_date
`

func (q *Queries) ListAccounts(ctx context.Context) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.AccountUuid,
			&i.Username,
			&i.Email,
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

const updateAccount = `-- name: UpdateAccount :one
UPDATE account 
   SET username = $1, 
       email = $2,
       updated_date = now()
 WHERE account_uuid = $3
 RETURNING account_uuid, username, email, created_date, updated_date, created_by, updated_by
`

type UpdateAccountParams struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccountUuid uuid.UUID `json:"account_uuid"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.Username, arg.Email, arg.AccountUuid)
	var i Account
	err := row.Scan(
		&i.AccountUuid,
		&i.Username,
		&i.Email,
		&i.CreatedDate,
		&i.UpdatedDate,
		&i.CreatedBy,
		&i.UpdatedBy,
	)
	return i, err
}
