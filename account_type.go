package saldo

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

// IsValid reports whether at is one of the five defined account types.
func (at AccountType) IsValid() bool {
	switch at {
	case Asset, Liability, Equity, Revenue, Expense:
		return true
	}
	return false
}
