package domain

import "context"

// LedgerRepository defines methods for ledger repository.
type LedgerRepository interface {
	Insert(context.Context, Ledger) error
	Get(context.Context, LedgerID) (Ledger, error)
	List(context.Context) ([]Ledger, error)
}

// AccountRepository defines methods for account repository.
type AccountRepository interface {
	Insert(context.Context, Account) error
	Get(context.Context, AccountID) (Account, error)
}
