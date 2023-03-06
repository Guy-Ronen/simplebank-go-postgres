package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Queries is a collection of all the SQL queries that we will use in our application.
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// function to execute a transaction
func (store *SQLStore) execTx(ctx context.Context, cbFn func(*Queries) error) error {

	// begin a transaction
	tx, err := store.db.BeginTx(ctx, nil)

	// if there is an error, return it
	if err != nil {
		return err
	}

	// create a new query object with the transaction
	query := New(tx)

	// execute the callback function
	err = cbFn(query)

	// if there is an error, rollback the transaction
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// if there is no error, commit the transaction
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
}

// function to transfer money between accounts from one account to the other.
// It creates a transfer record. add acount entry and update accounts balance within a singlr database transaction

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	// create a result object
	var result TransferTxResult

	// execute the transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create a transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// create a "from" transfer record
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create a "to" transfer record
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}


		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	AccountID1 int64,
	amount1 int64,
	AccountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     AccountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     AccountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
