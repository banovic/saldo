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

func ledgerFromDomain(l domain.Ledger) Ledger {
	return Ledger{
		LedgerID:           l.LedgerID,
		Name:               l.Name,
		FunctionalCurrency: l.FunctionalCurrency,
		ReportingTimeZone:  l.ReportingTimeZone,
		CreatedAt:          l.CreatedAt,
	}
}
