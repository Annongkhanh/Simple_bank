package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Mock db
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Postgres db
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}

}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Add new transfer record, add account entries and update accounts's balance
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update balance
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		argAccount1 := AddAccountBalanceParams{
			ID:     account1.ID,
			Amount: -arg.Amount,
		}

		// fmt.Println("account1 balance:", account1.Balance)

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		argAccount2 := AddAccountBalanceParams{
			ID:     account2.ID,
			Amount: arg.Amount,
		}

		if account1.ID < account2.ID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, argAccount1, argAccount2)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, argAccount2, argAccount1)
		}

		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{

		// })
		return nil
	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	argAccount1 AddAccountBalanceParams,
	argAccount2 AddAccountBalanceParams,
) (
	account1 Account,
	account2 Account,
	err error,
) {
	account1, err = q.AddAccountBalance(ctx, argAccount1)
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, argAccount2)
	if err != nil {
		return
	}
	return
}
