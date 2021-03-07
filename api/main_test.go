package api

import (
	"github.com/google/uuid"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/util"
)

var account db.Account = createRandomAccount()
var expectedAccount db.Account = db.Account{
	AccountUuid: uuid.New(),
	Username:    account.Username,
	Email:       account.Email,
}
var address db.Address = createRandomAddress()
var expectedAddress = db.Address{
	AddressUuid: uuid.New(),
	AccountUuid: account.AccountUuid,
	Name:        address.Name,
	Street:      address.Street,
	City:        address.City,
	State:       address.State,
	Zipcode:     address.Zipcode,
}

var accounts = createRandomAccountList(10)

func createRandomAccountList(size int64) []db.Account {
	result := make([]db.Account, size)
	for i := range result {
		result[i] = createRandomAccount()
	}
	return result
}

func createRandomAccount() db.Account {
	return db.Account{
		AccountUuid: uuid.New(),
		Username:    util.RandomUsername(),
		Email:       util.RandomEmail(),
	}
}

func createRandomAddress() db.Address {
	return db.Address{
		Name:    util.RandomString(10),
		Street:  util.RandomAlphaNumericString(30),
		City:    util.RandomString(10),
		State:   db.StateCA,
		Zipcode: util.RandomNumericString(5),
	}
}
