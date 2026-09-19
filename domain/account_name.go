package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	accountNameMinLen = 3
	accountNameMaxLen = 100
)

var (
	ErrInvalidAccountName = errors.New("invalid account name")
)

// AccountName is value object representing the name of the account.
type AccountName string

// NewAccountName constructs new account name from the given string.
func NewAccountName(str string) (AccountName, error) {
	an := AccountName(str)
	if err := an.Validate(); err != nil {
		return "", err
	}
	return an, nil
}

// Validate whether the account name is valid. Valid is if:
//   - non empty
//   - valid utf8 encoding
//   - does not have leading or trailing whitespaces
//   - between configured length bounds
//   - contains no control characters
func (an AccountName) Validate() error {
	str := string(an)
	if str == "" {
		return fmt.Errorf("%w: empty string", ErrInvalidAccountName)
	}
	if !utf8.ValidString(str) {
		return fmt.Errorf("%w: invalid utf8", ErrInvalidAccountName)
	}
	if str != strings.TrimSpace(str) {
		return fmt.Errorf("%w: leading or trailing whitespace", ErrInvalidAccountName)
	}
	rc := utf8.RuneCountInString(str)
	if rc < accountNameMinLen || rc > accountNameMaxLen {
		return fmt.Errorf("%w: min %d, max %d, got: %d", ErrInvalidAccountName, accountNameMinLen, accountNameMaxLen, rc)
	}
	for _, r := range str {
		if unicode.IsControl(r) {
			return fmt.Errorf("%w: invalid char (control): %U", ErrInvalidAccountName, r)
		}
	}
	return nil
}
