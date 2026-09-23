package domain

import (
	"errors"
	"fmt"
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
	if err := validateText(string(an), accountNameMinLen, accountNameMaxLen); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAccountName, err)
	}
	return nil
}
