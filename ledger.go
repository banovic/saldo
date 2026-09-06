package saldo

// LedgerID is ID of Ledger, it must be Valid and non-empty ("" is invalid).
type LedgerID string

// Ledger is a complete, self-contained set of books for one accounting entity.
type Ledger struct {
	LedgerID           LedgerID
	Name               string
	FunctionalCurrency Currency
}
