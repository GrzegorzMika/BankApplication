package db

import (
	"context"
	"math"
	"testing"

	"BankApplication/internal/util"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	_, accoount1 := createRandomAccount(t)
	_, accoount2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := util.RandomMoney()

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: accoount1.ID,
				ToAccountId:   accoount2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accoount1.ID, transfer.FromAccountID)
		require.Equal(t, accoount2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accoount1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accoount2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accoount1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accoount2.ID, toAccount.ID)

		diff1 := accoount1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - accoount2.Balance
		require.InDelta(t, diff1, diff2, 0.001)
		require.True(t, diff1 >= 0)

		k := int(math.Round(diff1 / amount))
		require.True(t, k >= 0 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), accoount1.ID)
	require.NoError(t, err)
	updatedAccount2, err := store.GetAccount(context.Background(), accoount2.ID)
	require.NoError(t, err)

	require.InDelta(t, accoount1.Balance-float64(n)*amount, updatedAccount1.Balance, 0.001)
	require.InDelta(t, accoount2.Balance+float64(n)*amount, updatedAccount2.Balance, 0.001)
}
