package domain

import (
	"errors"
	"fmt"
	"time"
	"uuid"
)

var (
	ErrNotEnoughPostings  = errors.New("not enough postings")
	ErrPostingsSumNotZero = errors.New("postings sum is not zero")
)

// JournalEntryID is ID of single Journal Entry, can be empty ("" is valid).
type JournalEntryID uuid.UUID

// JournalEntry is single accounting event - an atomic unit of Ledger.
//
// An entry consists of 2 or more Postings whose FunctionalAmounts sum to zero.
// It is written atomically all or none at all.
//
// Journal Entries are immutable. Delete or Update must not happen. An incorrect
// JournalEntry is corrected with new JournalEntry, and that new JournalEntry Reverses
// incorrect one (sets Reverses attribute).
type JournalEntry struct {
	JournalEntryID JournalEntryID
	LedgerID       LedgerID

	// IdempotencyKey is supplied by the client for each JournalEntry client wants to record.
	// It is used to guarantee that the JournalEntry is written exactly once.
	IdempotencyKey string

	// SourceDocumentReferenceID is reference to document which motivated JournalEntry (invoice, receipt number, etc.)
	SourceDocumentReferenceID string

	// When event happened in the world. Stored as UTC timestamp.
	OccurredAt time.Time

	// In which period it lands in. Stored as Date.
	PostedOn Date

	// When event was recorded by the system. Stored as UTC timestamp.
	RecordedAt time.Time

	// Description is natural text describing JournalEntry.
	Description string

	// Set when this JournalEntry reverses another JournalEntry, or empty otherwise.
	Reverses JournalEntryID

	// At least 2 Postings which must sum to 0.
	Postings []Posting
}

// Validate validates JournalEntry against Ledger's functional currency (fc).
// JournalEntry is valid:
// - functional currency (fc) is valid
// - at least 2 Postings
// - Postings sum to zero in their FunctionalAmount
func (je JournalEntry) Validate(fc Currency) error {
	if err := fc.Validate(); err != nil {
		return fmt.Errorf("functional currency: %w", err)
	}
	if len(je.Postings) < 2 {
		return fmt.Errorf("%w: %d posting(s)", ErrNotEnoughPostings, len(je.Postings))
	}
	ms := make([]Money, len(je.Postings))
	for i, p := range je.Postings {
		if p.FunctionalAmount.Currency != fc {
			return fmt.Errorf("posting %d: %w: %q (posting) vs %q (functional)", i, ErrCurrencyMismatch, p.FunctionalAmount.Currency, fc)
		}
		ms[i] = p.FunctionalAmount
	}
	sum, err := SumMoney(fc, ms)
	if err != nil {
		return fmt.Errorf("sum postings: %w", err)
	}
	if !sum.IsZeroAmount() {
		return fmt.Errorf("%w: %v", ErrPostingsSumNotZero, sum)
	}
	return nil
}
