package saldo

// PostingID is ID of single Posting, it must be Valid and non-empty ("" is invalid).
type PostingID string

// Posting is a line placed into an Account.
type Posting struct {
	PostingID      PostingID
	JournalEntryID JournalEntryID
	AccountID      AccountID

	// Amount in the Currency the event actually happened in.
	TransactionAmount Money

	// Amount in the Currency that will be booked - Ledger Functional Currency.
	FunctionalAmount Money

	// Exchange rate between Currencies for Document and Functional Amounts.
	ExchangeRateID ExchangeRateID
}
