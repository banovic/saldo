package domain

import (
	"errors"
	"fmt"
	"uuid"
)

var (
	ErrInvalidAccount = errors.New("invalid account")
)

// AccountID identifies a single Account.
type AccountID struct{ uuid.UUID }

// IsZero checks if account id is zero.
func (aid AccountID) IsZero() bool {
	return aid.UUID == uuid.Nil()
}

// Account is a named category that value flows through.
//
// There is no Balance stored in Account. Balance is calculated by
// summing Functional Amounts of every Posting that references Account.
type Account struct {
	AccountID AccountID
	LedgerID  LedgerID
	Type      AccountType
	Code      AccountCode
	Name      AccountName
}

// Validate returns error if account is not valid.
// LedgerID is assumed to belong to valid Ledger.
// Account is valid if:
//   - account id is not zero
//   - ledger id is not zero
//   - type is valid
//   - code is valid
//   - name is valid
func (a Account) Validate() error {
	if a.AccountID.IsZero() {
		return fmt.Errorf("%w: account id is zero", ErrInvalidAccount)
	}
	if a.LedgerID.IsZero() {
		return fmt.Errorf("%w: ledger id is zero", ErrInvalidAccount)
	}
	if err := a.Type.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAccount, err)
	}
	if err := a.Code.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAccount, err)
	}
	if err := a.Name.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAccount, err)
	}
	return nil
}
