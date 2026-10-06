package app

import (
	"context"

	"github.com/banovic/saldo/domain"
)

type fakeUnitOfWork struct {
	ledgers domain.LedgerRepository
}

func (f fakeUnitOfWork) Execute(ctx context.Context, work func(Repositories) error) error {
	return work(Repositories{Ledger: f.ledgers})
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
