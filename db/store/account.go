package store

import (
	"context"

	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// AccountStore interface with persistence operations for Account
type AccountStore interface {
	CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error)
	GetAccountByID(ctx context.Context, accountUUID uuid.UUID) (db.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (db.Account, error)
	ListAccounts(ctx context.Context) ([]db.Account, error)
	UpdateAccount(ctx context.Context, arg db.UpdateAccountParams) (db.Account, error)
}

// DbAccountStore implementation of account store executing the operations against a real SQL database
type DbAccountStore struct {
	queries db.Querier
}

// NewDbAccountStore returns a new db account store
func NewDbAccountStore(queries db.Querier) AccountStore {
	return &DbAccountStore{
		queries: queries,
	}
}

// CreateAccount create account
func (store *DbAccountStore) CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error) {
	return store.queries.CreateAccount(ctx, arg)
}

// GetAccountByID get account by id
func (store *DbAccountStore) GetAccountByID(ctx context.Context, accountUUID uuid.UUID) (db.Account, error) {
	return store.queries.GetAccountById(ctx, accountUUID)
}

// GetAccountByUsername get account by username
func (store *DbAccountStore) GetAccountByUsername(ctx context.Context, username string) (db.Account, error) {
	return store.queries.GetAccountByUsername(ctx, username)
}

// ListAccounts list accounts
func (store *DbAccountStore) ListAccounts(ctx context.Context) ([]db.Account, error) {
	return store.queries.ListAccounts(ctx)
}

// UpdateAccount update account
func (store *DbAccountStore) UpdateAccount(ctx context.Context, arg db.UpdateAccountParams) (db.Account, error) {
	return store.queries.UpdateAccount(ctx, arg)
}
