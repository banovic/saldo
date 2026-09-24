package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	accountCodeMinLen = 1
	accountCodeMaxLen = 100
)

var (
	ErrInvalidAccountCode = errors.New("invalid account code")
)

// AccountCode identifies account within Ledger's chart of accounts.
// Ledger's chart of accounts is set of all Accounts within Ledger.
// In database AccountCode is unique index on (LedgerID, AccountCode).
// Sr: Konto u Kontnom planu. Kontni plan je izveden od Kontnog okvira koji je zakonom propisan.
// Situation is similar in other jurisdictions - law defines template for chart of accounts
// and Ledger then implements that plan and produces its own chart of accounts (plan).
type AccountCode string

// Validate returns error if account code is not valid.
// TODO!!! - this will need more validation as things become more clearer about what accountants actually use.
// Account code is valid if:
//   - not empty / zero
func (ac AccountCode) Validate() error {
	str := string(ac)
	if str == "" {
		return fmt.Errorf("%w: empty string", ErrInvalidAccountCode)
	}
	if !utf8.ValidString(str) {
		return fmt.Errorf("%w: invalid utf8", ErrInvalidAccountCode)
	}
	if str != strings.TrimSpace(str) {
		return fmt.Errorf("%w: leading or trailing whitespace", ErrInvalidAccountCode)
	}
	rc := utf8.RuneCountInString(str)
	if rc < accountCodeMinLen || rc > accountCodeMaxLen {
		return fmt.Errorf("%w: min %d, max %d, got: %d", ErrInvalidAccountCode, accountCodeMinLen, accountCodeMaxLen, rc)
	}
	for _, r := range str {
		if unicode.IsControl(r) {
			return fmt.Errorf("%w: invalid char (control): %U", ErrInvalidAccountCode, r)
		}
	}
	return nil
}
