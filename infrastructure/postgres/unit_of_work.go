package postgres

import (
	"context"

	"github.com/banovic/saldo/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ app.UnitOfWork = (*UnitOfWork)(nil)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (uow *UnitOfWork) Execute(ctx context.Context, work func(app.Repositories) error) error {
	return pgx.BeginFunc(ctx, uow.pool, func(tx pgx.Tx) error {
		// Repositories are implemented and unexported in this package.
		// They are available to work function through app.Repositories interface.
		repositories := app.Repositories{Ledger: ledgerRepository{tx: tx}, Account: accountRepository{tx: tx}}
		return work(repositories)
	})
}
