package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines all functions to execute db queries and transactions
type Repository interface {
	ExtendedQuerier
	RunInTx(ctx context.Context, fn func(q ExtendedQuerier) error) error
}

// ExtendedQuerier extends any custom queries
type ExtendedQuerier interface {
	Querier
	CustomCreateUser(ctx context.Context, arg InsertUserParams) (User, error)
}

type PostgresRepository struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewPostgresRepository(connPool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

func (s *PostgresRepository) RunInTx(ctx context.Context, fn func(q ExtendedQuerier) error) error {

	tx, err := s.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(ctx); errRb != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return fmt.Errorf("error on rollback: %v, original error: %w", errRb, err)
		}
		return err
	}
	return tx.Commit(ctx)
}
