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
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
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
