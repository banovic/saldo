package domain

import (
	"errors"
	"fmt"
	"time"
	"uuid"
)

var (
	ErrInvalidJournalEntry = errors.New("invalid journal entry")
)

// JournalEntryID identifies a single Journal Entry.
type JournalEntryID struct{ uuid.UUID }

// IsZero checks if journal entry id is zero.
func (jeid JournalEntryID) IsZero() bool {
	return jeid.UUID == uuid.Nil()
}

// JournalEntry is single accounting event - an atomic unit of Ledger.
//
// An entry consists of 2 or more Postings whose FunctionalAmounts sum to zero.
// It is written atomically all or none at all.
//
// Journal Entries are immutable. Delete or Update must not happen.
// An incorrect JournalEntry is corrected with new JournalEntry, and that
// new JournalEntry Reverses incorrect one (sets Reverses attribute).
type JournalEntry struct {
	JournalEntryID JournalEntryID
	LedgerID       LedgerID

	// IdempotencyKey is supplied by the client for each JournalEntry client wants to record.
	// It is used to guarantee that the JournalEntry is written exactly once per ledger.
	// In database view IdempotencyKey is unique index on (LedgerID, IdempotencyKey).
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

// Validate validates JournalEntry against given Ledger.
// Ledger is assumed to be valid when passed in this method.
// JournalEntry is valid:
//   - JournalEntryID must not be zero value
//   - ledger id must match
//   - at least 2 Postings
//   - all postings must be valid
//   - Postings sum to zero in their FunctionalAmount
//   - Reverses must be different than JournalEntryID
//   - IdempotencyKey must not be empty
//   - PostedOn must be valid date
//   - PostedOn must not be before OccurredAt
func (je JournalEntry) Validate(l Ledger) error {
	if je.JournalEntryID.IsZero() {
		return fmt.Errorf("%w: journal entry id is zero", ErrInvalidJournalEntry)
	}
	if je.LedgerID != l.LedgerID {
		return fmt.Errorf("%w: mismatched ledger ids %v and %v", ErrInvalidJournalEntry, je.LedgerID, l.LedgerID)
	}
	if len(je.Postings) < 2 {
		return fmt.Errorf("%w: at least 2 postings required, got %d", ErrInvalidJournalEntry, len(je.Postings))
	}
	ms := make([]Money, len(je.Postings))
	for i, p := range je.Postings {
		if je.JournalEntryID != p.JournalEntryID {
			return fmt.Errorf("%w: posting %d: journal entry id %v and posting's journal entry id %v do not match", ErrInvalidJournalEntry, i, je.JournalEntryID, p.JournalEntryID)
		}
		if err := p.Validate(); err != nil {
			return fmt.Errorf("%w: posting %d: %w", ErrInvalidJournalEntry, i, err)
		}
		if p.FunctionalAmount.Currency != l.FunctionalCurrency {
			return fmt.Errorf("%w: posting %d: %w: %q (posting) vs %q (functional)", ErrInvalidJournalEntry, i, ErrCurrencyMismatch, p.FunctionalAmount.Currency, l.FunctionalCurrency)
		}
		ms[i] = p.FunctionalAmount
	}
	sum, err := SumMoney(l.FunctionalCurrency, ms)
	if err != nil {
		return fmt.Errorf("%w: postings sum: %w", ErrInvalidJournalEntry, err)
	}
	if !sum.IsZero() {
		return fmt.Errorf("%w: sum of postings functional amounts must be zero, got %v", ErrInvalidJournalEntry, sum)
	}
	if !je.Reverses.IsZero() && je.Reverses == je.JournalEntryID {
		return fmt.Errorf("%w: reverse (%v) is same as journal entry (%v)", ErrInvalidJournalEntry, je.Reverses, je.JournalEntryID)
	}
	if je.IdempotencyKey == "" {
		return fmt.Errorf("%w: idempotency key is required", ErrInvalidJournalEntry)
	}
	if err := je.PostedOn.Validate(); err != nil {
		return fmt.Errorf("%w: posted on: %w", ErrInvalidJournalEntry, err)
	}
	loc, err := l.ReportingTimeZone.Location()
	if err != nil {
		return fmt.Errorf("%w: ledger time zone: %w", ErrInvalidJournalEntry, err)
	}
	occurredOn, err := DateIn(je.OccurredAt, loc)
	if err != nil {
		return fmt.Errorf("%w: occurred on: %w", ErrInvalidJournalEntry, err)
	}
	if je.PostedOn.Before(occurredOn) {
		return fmt.Errorf("%w: posted on %v before occurred on %v", ErrInvalidJournalEntry, je.PostedOn, occurredOn)
	}
	return nil
}
