package db

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	// run n concurrent transfer transactions using a go routine
	n := 5
	amount := float64(10)

	// create a channel to receive errors and results from the go routines
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := range n {
		// used to pass context but not required
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			// to get the db transaction name
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errs <- err // return error to the go routine channel
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for range n {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), int32(fromEntry.ID))
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), int32(toEntry.ID))
		require.NoError(t, err)

		// check acc
		fromAcc := result.FromAccount
		require.NotEmpty(t, fromAcc)
		require.Equal(t, acc1.ID, fromAcc.ID)

		toAcc := result.ToAccount
		require.NotEmpty(t, toAcc)
		require.Equal(t, acc2.ID, toAcc.ID)

		// check acc balance
		diff1 := math.Round((acc1.Balance-fromAcc.Balance)*100) / 100 // round to two dp
		diff2 := math.Round((toAcc.Balance-acc2.Balance)*100) / 100
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, math.Mod(diff1, amount) == 0) // similar to diff1%amount == 0

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance-(float64(n)*amount), updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+(float64(n)*amount), updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	// run n concurrent transfer transactions using a go routine
	n := 10
	amount := float64(10)

	// create a channel to receive errors and results from the go routines
	errs := make(chan error)

	for i := range n {
		fromAccountId := acc1.ID
		toAccountId := acc2.ID

		if i%2 == 1 {
			fromAccountId = acc2.ID
			toAccountId = acc1.ID
		}

		go func() {
			// to get the db transaction name
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err // return error to the go routine channel
		}()
	}

	// check results
	for range n {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balances
	updatedAccount1, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance, updatedAccount2.Balance)
}
