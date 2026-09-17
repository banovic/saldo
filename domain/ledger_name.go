package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
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
	str := string(ln)
	if str == "" {
		return fmt.Errorf("%w: empty string", ErrInvalidLedgerName)
	}
	if !utf8.ValidString(str) {
		return fmt.Errorf("%w: invalid utf8", ErrInvalidLedgerName)
	}
	if str != strings.TrimSpace(str) {
		return fmt.Errorf("%w: leading or trailing whitespace", ErrInvalidLedgerName)
	}
	rc := utf8.RuneCountInString(str)
	if rc < ledgerNameMinLen || rc > ledgerNameMaxLen {
		return fmt.Errorf("%w: min %d, max %d, got: %d", ErrInvalidLedgerName, ledgerNameMinLen, ledgerNameMaxLen, rc)
	}
	for _, r := range str {
		if unicode.IsControl(r) {
			return fmt.Errorf("%w: invalid char (control): %U", ErrInvalidLedgerName, r)
		}
	}
	return nil
}
