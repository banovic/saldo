package app

import (
	"context"

	"github.com/banovic/saldo/domain"
)

type Repositories struct {
	Ledger  domain.LedgerRepository
	Account domain.AccountRepository
}

type UnitOfWork interface {
	Execute(ctx context.Context, work func(Repositories) error) error
}

func executeValue[T any](ctx context.Context, uow UnitOfWork, work func(Repositories) (T, error)) (T, error) {
	var result T

	err := uow.Execute(ctx, func(rs Repositories) error {
		var err error
		result, err = work(rs)
		return err
	})

	return result, err
}
