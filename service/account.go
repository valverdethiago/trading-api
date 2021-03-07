package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// AccountService service to handle business rules for accounts
type AccountService struct {
	queries db.Querier
}

// NewAccountService Creates new service for account
func NewAccountService(queries db.Querier) *AccountService {
	return &AccountService{
		queries: queries,
	}
}

// CreateAccount creates an account with address if provided
func (service *AccountService) CreateAccount(account db.Account, address *db.Address) (db.Account, db.Address, error) {
	var dbAccount db.Account
	var dbAddress db.Address

	if service.isUsernameAlreadyTaken(account.Username) {
		return dbAccount, dbAddress, errors.New("Username already taken")
	}
	arg := db.CreateAccountParams{
		Username: account.Username,
		Email:    account.Email,
	}
	dbAccount, err := service.queries.CreateAccount(context.Background(), arg)
	if err != nil {
		return dbAccount, dbAddress, err
	}
	if address != nil {
		dbAddress, err = service.createAddressForAccount(dbAccount, address)
	}
	return dbAccount, dbAddress, err
}

// ListAccounts list all available accounts
func (service *AccountService) ListAccounts() ([]db.Account, error) {
	return service.queries.ListAccounts(context.Background())
}

// GetAccountByID find account by id
func (service *AccountService) GetAccountByID(id uuid.UUID) (db.Account, error) {
	return service.queries.GetAccountById(context.Background(), id)
}

func (service *AccountService) isUsernameAlreadyTaken(Username string) bool {
	_, err := service.queries.GetAccountByUsername(context.Background(), Username)
	return err == nil || err != sql.ErrNoRows
}

// AssertAccountExists Returns the account with the given ID
func (service *AccountService) AssertAccountExists(ID uuid.UUID) (db.Account, error) {
	return service.queries.GetAccountById(context.Background(), ID)
}

func (service *AccountService) createAddressForAccount(account db.Account, address *db.Address) (db.Address, error) {
	var dbAddress db.Address
	addressArg := db.CreateAddressParams{
		Name:        address.Name,
		Street:      address.Street,
		City:        address.City,
		State:       address.State,
		Zipcode:     address.Zipcode,
		AccountUuid: account.AccountUuid,
	}
	dbAddress, err := service.queries.CreateAddress(context.Background(), addressArg)
	return dbAddress, err
}
