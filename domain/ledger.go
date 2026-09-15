package domain

import (
	"time"
	"uuid"
)

// LedgerID is ID of Ledger.
type LedgerID uuid.UUID

// Ledger is a complete, self-contained set of books for one accounting entity.
type Ledger struct {
	LedgerID           LedgerID
	Name               string
	FunctionalCurrency Currency

	// ReportingTimeZone is used to derive PostedOn dates in journal entries from OccurredAt timestamp.
	ReportingTimeZone TimeZone

	// When Ledger was created. Timestamp in UTC.
	CreatedAt time.Time
}
