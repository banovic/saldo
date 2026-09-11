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
