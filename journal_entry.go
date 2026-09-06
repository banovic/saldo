package saldo

import "time"

// JournalEntryID is ID of single Journal Entry, can be empty ("" is valid).
type JournalEntryID string

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
	// It is used to guarantee that the Transaction is written exactly once.
	IdempotencyKey string

	// SourceDocumentReferenceID is reference to document which motivated JournalEntry (invoice, receipt number, etc.)
	SourceDocumentReferenceID string

	// When event happened in the world.
	OccurredOn time.Time

	// In which period it lands in
	PostedOn time.Time

	// When event was recorded by the system.
	RecordedOn time.Time

	// Description is natural text describing JournalEntry.
	Description string

	// Set when this JournalEntry reverses another JournalEntry, or empty otherwise.
	Reverses JournalEntryID

	// At least 2 Postings which must sum to 0.
	Postings []Posting
}
