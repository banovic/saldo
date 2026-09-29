package app

import (
	"time"

	"github.com/banovic/saldo/domain"
)

// Ledger is dto (Data Transfer Object) that is sent back to caller.
type Ledger struct {
	LedgerID           domain.LedgerID
	Name               domain.LedgerName
	FunctionalCurrency domain.Currency
	ReportingTimeZone  domain.TimeZone
	CreatedAt          time.Time
}
