package domain

import "context"

// LedgerRepository defines methods for ledger repository.
type LedgerRepository interface {
	Insert(ctx context.Context, l Ledger) error
}

// AccountRepository defines methods for account repository.
type AccountRepository interface {
	Insert(ctx context.Context, a Account) error
}
