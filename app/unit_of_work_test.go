package app

import (
	"context"

	"github.com/banovic/saldo/domain"
)

type fakeUnitOfWork struct {
	ledgers  domain.LedgerRepository
	accounts domain.AccountRepository
}

func (f fakeUnitOfWork) Execute(ctx context.Context, work func(Repositories) error) error {
	return work(Repositories{Ledger: f.ledgers, Account: f.accounts})
}

// fakeLedgerRepository returns canned results; every method returns err when set.
type fakeLedgerRepository struct {
	ledger  domain.Ledger   // returned by Get
	ledgers []domain.Ledger // returned by List
	err     error
}

func (f fakeLedgerRepository) Insert(ctx context.Context, l domain.Ledger) error {
	return f.err
}

func (f fakeLedgerRepository) Get(ctx context.Context, id domain.LedgerID) (domain.Ledger, error) {
	if f.err != nil {
		return domain.Ledger{}, f.err
	}
	return f.ledger, nil
}

func (f fakeLedgerRepository) List(ctx context.Context) ([]domain.Ledger, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.ledgers, nil
}

// fakeAccountRepository returns canned results; every method returns err when set.
type fakeAccountRepository struct {
	account domain.Account // returned by Get
	err     error
}

func (f fakeAccountRepository) Insert(ctx context.Context, a domain.Account) error {
	return f.err
}

func (f fakeAccountRepository) Get(ctx context.Context, id domain.AccountID) (domain.Account, error) {
	if f.err != nil {
		return domain.Account{}, f.err
	}
	return f.account, nil
}
