package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// AddressService service to handle business rules for addresses
type AddressService struct {
	queries        db.Querier
	accountService *AccountService
}

// NewAddressService Creates new service for address
func NewAddressService(queries db.Querier, accountService *AccountService) *AddressService {
	return &AddressService{
		queries:        queries,
		accountService: accountService,
	}
}

// GetAddressByAccountID find account by id
func (service *AddressService) GetAddressByAccountID(ID uuid.UUID) (db.Address, error) {
	return service.queries.GetAddressByAccount(context.Background(), ID)
}

// CreateAddressForAccount creates an address for an account only if there's no address yet
func (service *AddressService) CreateAddressForAccount(ID uuid.UUID, address db.Address) (db.Address, error) {
	var dbAddress db.Address
	dbAccount, err := service.accountService.AssertAccountExists(ID)
	if err != nil {
		return dbAddress, err
	}
	if service.accountAlreadyHasAddress(ID) {
		return dbAddress, errors.New("Account has already an address")
	}

	arg := db.CreateAddressParams{
		Name:        address.Name,
		Street:      address.Street,
		City:        address.City,
		State:       address.State,
		Zipcode:     address.Zipcode,
		AccountUuid: dbAccount.AccountUuid,
	}
	return service.queries.CreateAddress(context.Background(), arg)
}

// UpdateAddressForAccount creates an address for an account only if there's no address yet
func (service *AddressService) UpdateAddressForAccount(ID uuid.UUID, address db.Address) (db.Address, error) {
	var dbAddress db.Address
	_, err := service.accountService.AssertAccountExists(ID)
	if err != nil {
		return dbAddress, err
	}
	dbAddress, err = service.GetAddressByAccountID(ID)
	if err != nil && err == sql.ErrNoRows {
		return dbAddress, err
	}

	arg := db.UpdateAddressParams{
		Name:        address.Name,
		Street:      address.Street,
		City:        address.City,
		State:       address.State,
		Zipcode:     address.Zipcode,
		AddressUuid: dbAddress.AddressUuid,
	}
	return service.queries.UpdateAddress(context.Background(), arg)
}

// getAddressByAccountID Returns the address attached to the account with the given ID
func (service *AddressService) getAddressByAccountID(ID uuid.UUID) (db.Address, error) {
	var dbAddress db.Address
	dbAccount, err := service.accountService.AssertAccountExists(ID)
	if err != nil {
		return dbAddress, err
	}
	return service.queries.GetAddressByAccount(context.Background(), dbAccount.AccountUuid)
}

func (service *AddressService) accountAlreadyHasAddress(ID uuid.UUID) bool {
	_, err := service.GetAddressByAccountID(ID)
	return err == nil || err != sql.ErrNoRows
}
