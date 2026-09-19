package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAccountType = errors.New("invalid account type")
)

// AccountType classifies an account by its role in the Accounting Equation:
// Assets = Liabilities + Equity
type AccountType string

// The five account types.
const (
	Asset     AccountType = "asset"
	Liability AccountType = "liability"
	Equity    AccountType = "equity"
	Revenue   AccountType = "revenue"
	Expense   AccountType = "expense"
)

// Validate returns error if AccountType is not valid.
// AccountType is valid:
//   - must be one of the five defined account types.
func (at AccountType) Validate() error {
	switch at {
	case Asset, Liability, Equity, Revenue, Expense:
		return nil
	}
	return fmt.Errorf("%w: %q", ErrInvalidAccountType, at)
}
