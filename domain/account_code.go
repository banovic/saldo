package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAccountCode = errors.New("invalid account code")
)

// AccountCode identifies account within Ledger's chart of accounts.
// Ledger's chart of accounts is set of all Accounts within Ledger.
// Sr: Konto u Kontnom planu. Kontni plan je izveden od Kontnog okvira koji je zakonom propisan.
// Situation is similar in other jurisdictions - law defines template for chart of accounts
// and Ledger then implements that plan and produces its own chart of accounts (plan).
type AccountCode string

// Validate returns error if account code is not valid.
// Account code is valid if:
//   - not empty / zero
func (ac AccountCode) Validate() error {
	if ac == "" {
		return fmt.Errorf("%w: must not be empty", ErrInvalidAccountCode)
	}
	// TODO!!! - this will need more validation as things become more clearer
	// about what accountants actually use.
	return nil
}
