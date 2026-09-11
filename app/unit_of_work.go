package app

import "github.com/banovic/saldo/domain"

type Repositories struct {
	Ledger  domain.LedgerRepository
	Account domain.AccountRepository
}

type UnitOfWork interface {
	Execute(work func(Repositories) error) error
}
