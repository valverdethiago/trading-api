package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valverdethiago/trading-api/util"
)

func createRandomAddress(t *testing.T) Address {
	account := createRandomAccount(t)
	arg := CreateAddressParams{
		Name:        util.RandomString(10),
		Street:      util.RandomString(20),
		City:        util.RandomString(10),
		State:       StateAK,
		Zipcode:     util.RandomZipcode(),
		AccountUuid: account.AccountUuid,
	}

	address, err := testQueries.CreateAddress(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotNil(t, address.AddressUuid)
	require.NotNil(t, address.CreatedDate)
	require.NotNil(t, address.UpdatedDate)
	require.Equal(t, arg.Name, address.Name)
	require.Equal(t, arg.Street, address.Street)
	require.Equal(t, arg.City, address.City)
	require.Equal(t, StateAK, address.State)
	require.Equal(t, arg.Zipcode, address.Zipcode)
	require.Equal(t, arg.AccountUuid, address.AccountUuid)
	return address
}

func TestCreateAddress(t *testing.T) {
	createRandomAddress(t)
}

func TestGetAddressById(t *testing.T) {
	address := createRandomAddress(t)
	dbAddress, err := testQueries.GetAddressById(context.Background(), address.AddressUuid)
	require.NoError(t, err)
	require.NotEmpty(t, dbAddress)
	require.Equal(t, dbAddress.Name, address.Name)
	require.Equal(t, dbAddress.Street, address.Street)
	require.Equal(t, dbAddress.City, address.City)
	require.Equal(t, dbAddress.State, address.State)
	require.Equal(t, dbAddress.Zipcode, address.Zipcode)
	require.Equal(t, dbAddress.AccountUuid, address.AccountUuid)
}

func TestUpdateAddress(t *testing.T) {
	address := createRandomAddress(t)
	arg := UpdateAddressParams{
		Name:        util.RandomString(10),
		Street:      util.RandomString(20),
		City:        util.RandomString(10),
		State:       StateAL,
		Zipcode:     util.RandomZipcode(),
		AddressUuid: address.AddressUuid,
	}

	dbAddress, err := testQueries.UpdateAddress(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dbAddress)
	require.NotEqual(t, dbAddress.Name, address.Name)
	require.NotEqual(t, dbAddress.Street, address.Street)
	require.NotEqual(t, dbAddress.City, address.City)
	require.NotEqual(t, dbAddress.State, address.State)
	require.NotEqual(t, dbAddress.Zipcode, address.Zipcode)
	require.Equal(t, dbAddress.AccountUuid, address.AccountUuid)
	require.NotEqual(t, address.UpdatedDate, dbAddress.UpdatedDate)
	require.NotEqual(t, dbAddress.CreatedDate, dbAddress.UpdatedDate)
}
func TestDeleteAddress(t *testing.T) {
	address := createRandomAddress(t)

	err := testQueries.DeleteAddressFromAccount(context.Background(), address.AccountUuid)
	require.NoError(t, err)
	dbAddress, err := testQueries.GetAddressByAccount(context.Background(), address.AccountUuid)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, dbAddress)
}
