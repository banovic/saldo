package domain

type LedgerRepository interface {
	Insert(l Ledger) error
}

type AccountRepository interface {
	Insert(a Account) error
}
