package domain

import (
	"errors"
	"fmt"
)

const (
	ledgerNameMinLen = 3
	ledgerNameMaxLen = 100
)

var (
	ErrInvalidLedgerName = errors.New("invalid ledger name")
)

// LedgerName is value object representing the name of ledger.
type LedgerName string

// NewLedgerName constructs new ledger name from given string.
func NewLedgerName(str string) (LedgerName, error) {
	ln := LedgerName(str)
	if err := ln.Validate(); err != nil {
		return "", err
	}
	return ln, nil
}

// Validate whether the ledger name is valid. Valid is if:
//   - non empty
//   - valid utf8 encoding
//   - does not have leading or trailing whitespaces
//   - between configured length bounds
//   - contains no control characters
func (ln LedgerName) Validate() error {
	if err := validateText(string(ln), ledgerNameMinLen, ledgerNameMaxLen); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidLedgerName, err)
	}
	return nil
}
