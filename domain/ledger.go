package domain

import (
	"errors"
	"fmt"
	"time"
	"uuid"
)

var (
	ErrInvalidLedger = errors.New("invalid ledger")
)

// LedgerID identifies a Ledger.
type LedgerID struct{ uuid.UUID }

// IsZero checks if ledger id is zero.
func (lid LedgerID) IsZero() bool {
	return lid.UUID == uuid.Nil()
}

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
// Ledger is valid if:
//   - ledger id is non-zero
//   - name is valid
//   - functional currency is valid
//   - reporting time zone is valid
func (l Ledger) Validate() error {
	if l.LedgerID.IsZero() {
		return fmt.Errorf("%w: ledger id is zero", ErrInvalidLedger)
	}

	if err := l.Name.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidLedger, err)
	}

	if err := l.FunctionalCurrency.Validate(); err != nil {
		return fmt.Errorf("functional currency: %w", err)
	}

	if err := l.ReportingTimeZone.Validate(); err != nil {
		return fmt.Errorf("reporting time zone: %w", err)
	}
	return nil
}
