package store

import (
	"context"

	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// AddressStore interface with persistence operations for Account
type AddressStore interface {
	CreateAddress(ctx context.Context, arg db.CreateAddressParams) (db.Address, error)
	DeleteAddressFromAccount(ctx context.Context, accountUUID uuid.UUID) error
	GetAddressByAccount(ctx context.Context, accountUUID uuid.UUID) (db.Address, error)
	GetAddressByID(ctx context.Context, accountUUID uuid.UUID) (db.Address, error)
	UpdateAddress(ctx context.Context, arg db.UpdateAddressParams) (db.Address, error)
}

// DbAddressStore implementation of account store executing the operations against a real SQL database
type DbAddressStore struct {
	queries db.Querier
}

// NewDbAddressStore builds a new instance of db account store
func NewDbAddressStore(queries db.Querier) AddressStore {
	return &DbAddressStore{
		queries: queries,
	}
}

// CreateAddress creates an address
func (addressStore *DbAddressStore) CreateAddress(ctx context.Context, arg db.CreateAddressParams) (db.Address, error) {
	return addressStore.queries.CreateAddress(ctx, arg)
}

// DeleteAddressFromAccount removes an address from an account
func (addressStore *DbAddressStore) DeleteAddressFromAccount(ctx context.Context, accountUUID uuid.UUID) error {
	return addressStore.queries.DeleteAddressFromAccount(ctx, accountUUID)
}

// GetAddressByAccount returns an address attached to the account
func (addressStore *DbAddressStore) GetAddressByAccount(ctx context.Context, accountUUID uuid.UUID) (db.Address, error) {
	return addressStore.queries.GetAddressByAccount(ctx, accountUUID)
}

// GetAddressByID gets and address by its id
func (addressStore *DbAddressStore) GetAddressByID(ctx context.Context, accountUUID uuid.UUID) (db.Address, error) {
	return addressStore.queries.GetAddressById(ctx, accountUUID)
}

// UpdateAddress updates an address
func (addressStore *DbAddressStore) UpdateAddress(ctx context.Context, arg db.UpdateAddressParams) (db.Address, error) {
	return addressStore.queries.UpdateAddress(ctx, arg)
}
