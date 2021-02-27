// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createAccountStmt, err = db.PrepareContext(ctx, createAccount); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAccount: %w", err)
	}
	if q.createAddressStmt, err = db.PrepareContext(ctx, createAddress); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAddress: %w", err)
	}
	if q.createTradeStmt, err = db.PrepareContext(ctx, createTrade); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTrade: %w", err)
	}
	if q.deleteAddressFromAccountStmt, err = db.PrepareContext(ctx, deleteAddressFromAccount); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAddressFromAccount: %w", err)
	}
	if q.getAccountByIdStmt, err = db.PrepareContext(ctx, getAccountById); err != nil {
		return nil, fmt.Errorf("error preparing query GetAccountById: %w", err)
	}
	if q.getAddressByAccountStmt, err = db.PrepareContext(ctx, getAddressByAccount); err != nil {
		return nil, fmt.Errorf("error preparing query GetAddressByAccount: %w", err)
	}
	if q.getAddressByIdStmt, err = db.PrepareContext(ctx, getAddressById); err != nil {
		return nil, fmt.Errorf("error preparing query GetAddressById: %w", err)
	}
	if q.getTradeByIdStmt, err = db.PrepareContext(ctx, getTradeById); err != nil {
		return nil, fmt.Errorf("error preparing query GetTradeById: %w", err)
	}
	if q.listAccountsStmt, err = db.PrepareContext(ctx, listAccounts); err != nil {
		return nil, fmt.Errorf("error preparing query ListAccounts: %w", err)
	}
	if q.listTradesByAccountStmt, err = db.PrepareContext(ctx, listTradesByAccount); err != nil {
		return nil, fmt.Errorf("error preparing query ListTradesByAccount: %w", err)
	}
	if q.updateAccountStmt, err = db.PrepareContext(ctx, updateAccount); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAccount: %w", err)
	}
	if q.updateAddressStmt, err = db.PrepareContext(ctx, updateAddress); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAddress: %w", err)
	}
	if q.updateTradeStmt, err = db.PrepareContext(ctx, updateTrade); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTrade: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createAccountStmt != nil {
		if cerr := q.createAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAccountStmt: %w", cerr)
		}
	}
	if q.createAddressStmt != nil {
		if cerr := q.createAddressStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAddressStmt: %w", cerr)
		}
	}
	if q.createTradeStmt != nil {
		if cerr := q.createTradeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTradeStmt: %w", cerr)
		}
	}
	if q.deleteAddressFromAccountStmt != nil {
		if cerr := q.deleteAddressFromAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAddressFromAccountStmt: %w", cerr)
		}
	}
	if q.getAccountByIdStmt != nil {
		if cerr := q.getAccountByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountByIdStmt: %w", cerr)
		}
	}
	if q.getAddressByAccountStmt != nil {
		if cerr := q.getAddressByAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAddressByAccountStmt: %w", cerr)
		}
	}
	if q.getAddressByIdStmt != nil {
		if cerr := q.getAddressByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAddressByIdStmt: %w", cerr)
		}
	}
	if q.getTradeByIdStmt != nil {
		if cerr := q.getTradeByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTradeByIdStmt: %w", cerr)
		}
	}
	if q.listAccountsStmt != nil {
		if cerr := q.listAccountsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAccountsStmt: %w", cerr)
		}
	}
	if q.listTradesByAccountStmt != nil {
		if cerr := q.listTradesByAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listTradesByAccountStmt: %w", cerr)
		}
	}
	if q.updateAccountStmt != nil {
		if cerr := q.updateAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAccountStmt: %w", cerr)
		}
	}
	if q.updateAddressStmt != nil {
		if cerr := q.updateAddressStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAddressStmt: %w", cerr)
		}
	}
	if q.updateTradeStmt != nil {
		if cerr := q.updateTradeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTradeStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                           DBTX
	tx                           *sql.Tx
	createAccountStmt            *sql.Stmt
	createAddressStmt            *sql.Stmt
	createTradeStmt              *sql.Stmt
	deleteAddressFromAccountStmt *sql.Stmt
	getAccountByIdStmt           *sql.Stmt
	getAddressByAccountStmt      *sql.Stmt
	getAddressByIdStmt           *sql.Stmt
	getTradeByIdStmt             *sql.Stmt
	listAccountsStmt             *sql.Stmt
	listTradesByAccountStmt      *sql.Stmt
	updateAccountStmt            *sql.Stmt
	updateAddressStmt            *sql.Stmt
	updateTradeStmt              *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                           tx,
		tx:                           tx,
		createAccountStmt:            q.createAccountStmt,
		createAddressStmt:            q.createAddressStmt,
		createTradeStmt:              q.createTradeStmt,
		deleteAddressFromAccountStmt: q.deleteAddressFromAccountStmt,
		getAccountByIdStmt:           q.getAccountByIdStmt,
		getAddressByAccountStmt:      q.getAddressByAccountStmt,
		getAddressByIdStmt:           q.getAddressByIdStmt,
		getTradeByIdStmt:             q.getTradeByIdStmt,
		listAccountsStmt:             q.listAccountsStmt,
		listTradesByAccountStmt:      q.listTradesByAccountStmt,
		updateAccountStmt:            q.updateAccountStmt,
		updateAddressStmt:            q.updateAddressStmt,
		updateTradeStmt:              q.updateTradeStmt,
	}
}
