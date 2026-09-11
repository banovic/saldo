package domain

import "context"

type LedgerRepository interface {
	Insert(ctx context.Context, l Ledger) error
}

type AccountRepository interface {
	Insert(ctx context.Context, a Account) error
}
