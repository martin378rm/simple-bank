//package db
//
//import (
//	"context"
//	"fmt"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//)
//
//func TestTransferTx2(t *testing.T) {
//	store := NewStore(testDB)
//
//	account1 := createRandomAccount(t)
//	account2 := createRandomAccount(t)
//
//	fmt.Printf(">> before : from %v --> to %v\n", account1.Balance, account2.Balance)
//
//	// run conccurent transfer transactions
//	n := 2
//	amount := int64(10)
//
//	errs := make(chan error)
//	results := make(chan TransferTxResult)
//	for i := 0; i < n; i++ {
//		txName := fmt.Sprintf("tx %d", i+1)
//		go func() {
//			ctx := context.WithValue(context.Background(), txKey, txName)
//			result, err := store.TransferTx(ctx, TransferTxParams{
//				FromAccountID: account1.ID,
//				ToAccountID:   account2.ID,
//				Amount:        amount,
//			})
//			errs <- err
//			results <- result
//		}()
//	}
//
//	// check results
//	existed := make(map[int]bool)
//	for i := 0; i < n; i++ {
//		err := <-errs
//		require.NoError(t, err)
//
//		result := <-results
//		require.NotEmpty(t, result)
//
//		// check transfer
//		transfer := result.Transfer
//		require.NotEmpty(t, transfer)
//		require.Equal(t, account1.ID, transfer.FromAccountID)
//		require.Equal(t, account2.ID, transfer.ToAccountID)
//		require.Equal(t, amount, transfer.Amount)
//		require.NotZero(t, transfer.ID)
//		require.NotZero(t, transfer.CreatedAt)
//
//		_, err = store.GetTransfer(context.Background(), transfer.ID)
//		require.NoError(t, err)
//
//		// check entries
//		fromEntry := result.FromEntry
//		require.NotEmpty(t, fromEntry)
//		require.Equal(t, account1.ID, fromEntry.AccountID)
//		require.Equal(t, -amount, fromEntry.Amount)
//		require.NotZero(t, fromEntry.AccountID)
//		require.NotZero(t, fromEntry.CreatedAt)
//
//		_, err = store.GetEntry(context.Background(), fromEntry.ID)
//		require.NoError(t, err)
//
//		toEntry := result.ToEntry
//		require.NotEmpty(t, toEntry)
//		require.Equal(t, account2.ID, toEntry.AccountID)
//		require.Equal(t, amount, toEntry.Amount)
//		require.NotZero(t, toEntry.AccountID)
//		require.NotZero(t, toEntry.CreatedAt)
//
//		_, err = store.GetEntry(context.Background(), toEntry.ID)
//		require.NoError(t, err)
//
//		// check accounts
//		fromAccount := result.FromAccount
//		require.NotEmpty(t, fromAccount)
//		require.Equal(t, account1.ID, fromAccount.ID)
//
//		toAccount := result.ToAccount
//		require.NotEmpty(t, toAccount)
//		require.Equal(t, account2.ID, toAccount.ID)
//
//		// check account balance
//		fmt.Printf(">> tx : from %v --> to %v\n", fromAccount.Balance, toAccount.Balance)
//
//		diff1 := account1.Balance - fromAccount.Balance
//		diff2 := toAccount.Balance - account2.Balance
//		require.Equal(t, diff1, diff2)
//		require.True(t, diff1 > 0)
//		require.True(t, diff1%amount == 0)
//
//		k := int(diff1 / amount)
//		require.True(t, k >= 1 && k <= n)
//		require.NotContains(t, existed, k)
//		existed[k] = true
//
//		// check the final updated balance
//		updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
//		require.NoError(t, err)
//
//		updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
//		require.NoError(t, err)
//
//		fmt.Printf(">> after : from %v --> to %v\n", updateAccount1.Balance, updateAccount2.Balance)
//
//		require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
//		require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
//		// TODO check balance
//		// 7:48
//
//	}
//
//}

package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTxx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
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

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", result.FromAccount.Balance, result.ToAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
