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

func ledgersFromDomain(ls []domain.Ledger) []Ledger {
	d := make([]Ledger, 0, len(ls))
	for _, l := range ls {
		d = append(d, ledgerFromDomain(l))
	}
	return d
}

// Account DTO.
type Account struct {
	AccountID domain.AccountID
	LedgerID  domain.LedgerID
	Type      domain.AccountType
	Code      domain.AccountCode
	Name      domain.AccountName
	CreatedAt time.Time
}

func accountFromDomain(a domain.Account) Account {
	return Account{
		AccountID: a.AccountID,
		LedgerID:  a.LedgerID,
		Type:      a.Type,
		Code:      a.Code,
		Name:      a.Name,
		CreatedAt: a.CreatedAt,
	}
}

func accountsFromDomain(as []domain.Account) []Account {
	d := make([]Account, 0, len(as))
	for _, a := range as {
		d = append(d, accountFromDomain(a))
	}
	return d
}
