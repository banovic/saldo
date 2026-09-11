package domain

// AccountID is ID of Account, it must be Valid and non-empty ("" is invalid).
type AccountID string

// IsValid checks if AccountID is valid.
// AccountID is valid:
// - must be non-empty ("" is invalid).
func (aid AccountID) IsValid() bool {
	return aid != ""
}

// Account is a named category that value flows through, identified
// within a Ledger and classified as one of the five Types.
//
// There is no Balance stored in Account. Balance is calculated by
// summing Functional Amounts of every Posting that references Account.
type Account struct {
	AccountID AccountID
	LedgerID  LedgerID
	Type      AccountType
	Code      string // Optional, value from national regulations etc.
	Name      string
}
