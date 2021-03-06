package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valverdethiago/trading-api/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotNil(t, account.AccountUuid)
	require.NotNil(t, account.CreatedDate)
	require.NotNil(t, account.UpdatedDate)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.Username, account.Username)

	return account
}

func findAccountInList(accounts []Account, account Account) Account {
	var result Account
	for _, element := range accounts {
		if element.AccountUuid == account.AccountUuid {
			result = element
		}
	}
	return result
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	dbAccount, err := testQueries.GetAccountById(context.Background(), account.AccountUuid)
	require.NoError(t, err)
	require.NotEmpty(t, dbAccount)
	require.Equal(t, account.AccountUuid, dbAccount.AccountUuid)
	require.Equal(t, account.CreatedDate, account.CreatedDate)
	require.Equal(t, account.Email, dbAccount.Email)
	require.Equal(t, account.Username, dbAccount.Username)
	require.Equal(t, account.UpdatedDate, dbAccount.UpdatedDate)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	arg := UpdateAccountParams{
		Username:    util.RandomUsername(),
		Email:       util.RandomEmail(),
		AccountUuid: account.AccountUuid,
	}

	dbAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dbAccount)
	require.NotEqual(t, account.Username, dbAccount.Username)
	require.NotEqual(t, account.Email, dbAccount.Email)
	require.NotEqual(t, dbAccount.CreatedDate, dbAccount.UpdatedDate)
}

func TestListAccounts(t *testing.T) {
	var accounts [10]Account
	for i := 0; i >= len(accounts); i++ {
		accounts[i] = createRandomAccount(t)
	}
	dbAccounts, err := testQueries.ListAccounts(context.Background())
	require.NoError(t, err)
	for _, account := range accounts {
		dbAccount := findAccountInList(dbAccounts, account)
		assert.Equal(t, dbAccount, account)
	}
}
