package domain

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrInvalidLedgerName = errors.New("invalid ledger name")
)

// LedgerName is value object representing name of ledger.
type LedgerName string

func (ln LedgerName) Validate() error {
	// TODO: must be utf8 encoded; trim whitespaces; maybe limit to certain chars?
	if utf8.RuneCountInString(string(ln)) < 3 || utf8.RuneCountInString(string(ln)) > 100 {
		return ErrInvalidLedgerName
	}
	return nil
}
