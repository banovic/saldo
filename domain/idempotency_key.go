package domain

import (
	"errors"
	"fmt"
)

const (
	idempotencyKeyMinLen = 3
	idempotencyKeyMaxLen = 100
)

var (
	ErrInvalidIdempotencyKey = errors.New("invalid idempotency key")
)

// IdempotencyKey is value object representing the idempotecy key created and sent by caller.
// It is used to guarantee that a JournalEntry is written exactly once per Ledger.
// IdempotencyKey is ASCII encoded string.
type IdempotencyKey string

// NewIdempotencyKey constructs new idempotency key from given string.
func NewIdempotencyKey(str string) (IdempotencyKey, error) {
	ik := IdempotencyKey(str)
	if err := ik.Validate(); err != nil {
		return "", err
	}
	return ik, nil
}

// Validate whether the idempotency key is valid. Valid is if:
//   - between configured length bounds
//   - contains only chars: a..zA..Z0..9_-
func (ik IdempotencyKey) Validate() error {
	l := len(ik)
	if l < idempotencyKeyMinLen || l > idempotencyKeyMaxLen {
		return fmt.Errorf("%w: min %d, max %d, got: %d", ErrInvalidIdempotencyKey, idempotencyKeyMinLen, idempotencyKeyMaxLen, l)
	}
	for i, c := range ik {
		switch {
		case c >= 'a' && c <= 'z': // allow
		case c >= 'A' && c <= 'Z': // allow
		case c >= '0' && c <= '9': // allow
		case c == '_' || c == '-': // allow
		default:
			return fmt.Errorf("%w: invalid char %q at index %d", ErrInvalidIdempotencyKey, c, i)
		}
	}
	return nil
}
