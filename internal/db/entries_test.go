package db

import (
	"context"
	"testing"

	"BankApplication/internal/util"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (CreateEntryParams, Entry) {
	_, user := createRandomUser(t)
	arg := CreateEntryParams{
		AccountID: util.RandomAccountID(),
		Amount:    util.RandomMoney(),
	}

	// due to the foreign key constraint on accounts, we need to create the account first
	query := `INSERT INTO accounts (id, owner, balance, currency) VALUES ($1, $2, $3, $4)`
	_, err := testQueries.db.Exec(context.Background(), query, arg.AccountID,
		user.Username, util.RandomMoney(), util.RandomCurrency())
	require.NoError(t, err)

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	return arg, entry
}

func TestCreateEntry(t *testing.T) {
	for i := 0; i < 100; i++ {
		arg, entry := createRandomEntry(t)
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.Equal(t, arg.Amount, entry.Amount)
		require.GreaterOrEqual(t, entry.Amount, float64(0))
		require.NotZero(t, entry.ID)
	}
}

func TestGetEntry(t *testing.T) {
	for i := 0; i < 100; i++ {
		arg, createdEntry := createRandomEntry(t)
		entry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.Equal(t, arg.Amount, entry.Amount)
		require.GreaterOrEqual(t, entry.Amount, float64(0))
		require.NotZero(t, entry.ID)
	}
}

func TestListEntries(t *testing.T) {
	_, user := createRandomUser(t)
	accountID := util.RandomAccountID()

	// due to the foreign key constraint on accounts, we need to create the account first
	query := `INSERT INTO accounts (id, owner, balance, currency) VALUES ($1, $2, $3, $4)`
	_, err := testQueries.db.Exec(context.Background(), query, accountID,
		user.Username, util.RandomMoney(), util.RandomCurrency())
	require.NoError(t, err)

	for i := 0; i < 20; i++ {
		arg := CreateEntryParams{
			AccountID: accountID,
			Amount:    util.RandomMoney(),
		}
		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}
	arg := ListEntriesParams{
		AccountID: accountID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
