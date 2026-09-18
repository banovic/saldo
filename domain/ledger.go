package domain

import (
	"fmt"
	"time"
	"uuid"
)

// LedgerID identifies a Ledger.
type LedgerID struct{ uuid.UUID }

// Ledger is a complete, self-contained set of books for one accounting entity.
type Ledger struct {
	LedgerID           LedgerID
	Name               LedgerName
	FunctionalCurrency Currency

	// ReportingTimeZone is used to derive PostedOn dates in journal entries from OccurredAt timestamp.
	ReportingTimeZone TimeZone

	// CreatedAt is timestamp in UTC when the Ledger was created (app time is source of truth).
	CreatedAt time.Time
}

// Validate returns error if ledger is not valid.
func (l Ledger) Validate() error {
	if err := l.Name.Validate(); err != nil {
		return err
	}

	if err := l.FunctionalCurrency.Validate(); err != nil {
		return fmt.Errorf("functional currency: %w", err)
	}

	if err := l.ReportingTimeZone.Validate(); err != nil {
		return fmt.Errorf("reporting time zone: %w", err)
	}
	return nil
}
