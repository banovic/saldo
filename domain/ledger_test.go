package domain

import (
	"errors"
	"testing"
	"time"
	"uuid"
)

func TestLedgerValidate(t *testing.T) {
	valid := Ledger{
		LedgerID:           LedgerID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")},
		Name:               "Acme",
		FunctionalCurrency: USD,
		ReportingTimeZone:  "Europe/Belgrade",
		CreatedAt:          time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC),
	}
	with := func(f func(*Ledger)) Ledger {
		l := valid
		f(&l)
		return l
	}

	testCases := []struct {
		name    string
		ledger  Ledger
		wantErr error
	}{
		{"valid", valid, nil},
		{"invalid name", with(func(l *Ledger) { l.Name = "ab" }), ErrInvalidLedgerName},
		{"invalid functional currency", with(func(l *Ledger) { l.FunctionalCurrency = "XXX" }), ErrInvalidCurrency},
		{"empty reporting time zone", with(func(l *Ledger) { l.ReportingTimeZone = "" }), ErrInvalidTimeZone},
		{"unknown reporting time zone", with(func(l *Ledger) { l.ReportingTimeZone = "Europe/Nowhere" }), ErrInvalidTimeZone},
		{"name is checked before currency", with(func(l *Ledger) {
			l.Name = ""
			l.FunctionalCurrency = "XXX"
		}), ErrInvalidLedgerName},
		{"currency is checked before time zone", with(func(l *Ledger) {
			l.FunctionalCurrency = "XXX"
			l.ReportingTimeZone = ""
		}), ErrInvalidCurrency},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.ledger.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%+v.Validate() = %v, want %v", tc.ledger, err, tc.wantErr)
			}
		})
	}
}
