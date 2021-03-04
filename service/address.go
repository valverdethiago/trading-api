package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// AddressService service to handle business rules for addresses
type AddressService struct {
	queries *db.Queries
}

// NewAddressService Creates new service for address
func NewAddressService(queries *db.Queries) *AddressService {
	return &AddressService{
		queries: queries,
	}
}

// GetAddressByAccountID find account by id
func (service *AddressService) GetAddressByAccountID(ID string) (db.Address, error) {
	var dbAddress db.Address
	log.Printf("Trying to parse id %s", ID)
	uuid, err := uuid.Parse(ID)
	if err != nil {
		return dbAddress, errors.New("Invalid ID")
	}
	log.Printf("trying to fetch address for account with id %s", uuid)
	return service.queries.GetAddressByAccount(context.Background(), uuid)
}

// CreateAddressForAccount creates an address for an account only if there's no address yet
func (service *AddressService) CreateAddressForAccount(ID string, address db.Address) (db.Address, error) {
	var dbAddress db.Address
	dbAccount, err := service.assertAccountExists(ID)
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
func (service *AddressService) UpdateAddressForAccount(ID string, address db.Address) (db.Address, error) {
	var dbAddress db.Address
	_, err := service.assertAccountExists(ID)
	if err != nil {
		return dbAddress, err
	}
	dbAddress, err = service.GetAddressByAccountID(ID)
	if err != nil && err == sql.ErrNoRows {
		return dbAddress, errors.New("Account doesn't have an address yet")
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
func (service *AddressService) getAddressByAccountID(ID string) (db.Address, error) {
	var dbAddress db.Address
	dbAccount, err := service.assertAccountExists(ID)
	if err != nil {
		return dbAddress, err
	}
	return service.queries.GetAddressByAccount(context.Background(), dbAccount.AccountUuid)
}

func parseUUID(ID string) (uuid.UUID, error) {
	log.Printf("Trying to parse id %s", ID)
	var result uuid.UUID
	result, err := uuid.Parse(ID)
	if err != nil {
		return result, errors.New("Invalid ID")
	}
	return result, nil
}

func (service *AddressService) assertAccountExists(ID string) (db.Account, error) {
	var dbAccount db.Account
	uuid, err := parseUUID(ID)
	if err != nil {
		return dbAccount, err
	}
	dbAccount, err = service.queries.GetAccountById(context.Background(), uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return dbAccount, errors.New("No account found for the given id")
		}
		return dbAccount, err
	}
	log.Printf("Found account with id %s", uuid)
	return dbAccount, nil
}

func (service *AddressService) accountAlreadyHasAddress(ID string) bool {
	_, err := service.GetAddressByAccountID(ID)
	return err == nil || err != sql.ErrNoRows
}
