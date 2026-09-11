package postgres

import (
	"context"
	"errors"

	"github.com/banovic/saldo/app"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedToBeginTransaction = errors.New("failed to start transaction")
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (uow *UnitOfWork) Execute(ctx context.Context, work func(app.Repositories) error) error {
	tx, err := uow.pool.Begin(ctx)
	if err != nil {
		return ErrFailedToBeginTransaction
	}

	// If Commit() is successfully called first, Rollback() does nothing.
	defer tx.Rollback(ctx)

	// Create repositories with tx.
	repositories := app.Repositories{Ledger: ledgerRepository{db: &tx}, Account: accountRepository{db: &tx}}

	if err := work(repositories); err != nil {
		return err
	}

	// Commit changes made by repositories.
	return tx.Commit(ctx)
}
