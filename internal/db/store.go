package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is the interface for db operations
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// PostgresSQLStore provides all function to execute db queries against PostgresSQL db and transactions.
type PostgresSQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new Store instance
func NewStore(db *pgxpool.Pool) Store {
	return &PostgresSQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes the given function in a transaction.
func (s *PostgresSQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Join(err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountId int64   `json:"from_account_id"`
	ToAccountId   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx transfers the funds from the sender account to the receiver account.
// It creates a new transfer, accounts entries and updates the account balances within a single transaction
func (s *PostgresSQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(queries *Queries) error {
		var err error
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, queries, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, queries, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, queries *Queries, fromAccountId int64, fromAmount float64, toAccountId int64, toAmount float64) (account1 Account, account2 Account, err error) {
	account1, err = queries.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     fromAccountId,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}

	account2, err = queries.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     toAccountId,
		Amount: toAmount,
	})
	return
}
