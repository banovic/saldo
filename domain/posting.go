package domain

import (
	"errors"
	"fmt"
	"uuid"
)

var (
	ErrInvalidPosting = errors.New("invalid posting")
)

// PostingID identifies a posting.
type PostingID struct{ uuid.UUID }

// Posting is a line placed into an Account.
type Posting struct {
	PostingID      PostingID
	JournalEntryID JournalEntryID
	AccountID      AccountID

	// Amount in the Currency the event actually happened in.
	TransactionAmount Money

	// Amount in the Currency that will be booked - Ledger Functional Currency.
	FunctionalAmount Money

	// Exchange rate between Currencies for Transaction and Functional Amount's Currencies.
	ExchangeRateID ExchangeRateID
}

// IsDebit checks if amount of functional money is greater than zero.
func (p Posting) IsDebit() bool {
	return p.FunctionalAmount.MinorUnits > 0
}

// IsCredit checks if amount of functional money is less than zero.
func (p Posting) IsCredit() bool {
	return p.FunctionalAmount.MinorUnits < 0
}

// Validate returns error if posting is not valid.
// Posting is valid if:
//   - transactional amount is valid
//   - functional amount is valid
//   - both amounts have same sign
//   - if exchange rate id is zero, amounts are identical (same amount and currency)
//   - if exchange rate id is non-zero, the amounts are in different currencies
//   - if the amounts have same currency, they are identical (same minor units)
func (p Posting) Validate() error {
	if err := p.TransactionAmount.Validate(); err != nil {
		return fmt.Errorf("%w: transactional amount: %w", ErrInvalidPosting, err)
	}
	if err := p.FunctionalAmount.Validate(); err != nil {
		return fmt.Errorf("%w: functional amount: %w", ErrInvalidPosting, err)
	}
	if p.FunctionalAmount.Sign() != p.TransactionAmount.Sign() {
		return fmt.Errorf("%w: transaction amount sign %v does not match functional amount sign %v", ErrInvalidPosting, p.TransactionAmount, p.FunctionalAmount)
	}
	if p.ExchangeRateID.IsZero() && p.FunctionalAmount != p.TransactionAmount {
		return fmt.Errorf("%w: no exchange rate, amounts are different %v and %v", ErrInvalidPosting, p.FunctionalAmount, p.TransactionAmount)
	}
	if !p.ExchangeRateID.IsZero() && p.FunctionalAmount.Currency == p.TransactionAmount.Currency {
		return fmt.Errorf("%w: exchange rate set, both amounts are in same currency: %q", ErrInvalidPosting, p.FunctionalAmount.Currency)
	}
	if p.FunctionalAmount.Currency == p.TransactionAmount.Currency && p.FunctionalAmount != p.TransactionAmount {
		return fmt.Errorf("%w: same currency, amounts differ: %v and %v", ErrInvalidPosting, p.FunctionalAmount, p.TransactionAmount)
	}
	return nil
}
