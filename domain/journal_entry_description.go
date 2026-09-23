package domain

import (
	"errors"
	"fmt"
)

const (
	journalEntryDescriptionMinLen = 3
	journalEntryDescriptionMaxLen = 100
)

var (
	ErrInvalidJournalEntryDescription = errors.New("invalid journal entry description")
)

// JournalEntryDescription is value object representing the description of journal entry.
type JournalEntryDescription string

// NewJournalEntryDescription constructs new journal entry description from given string.
func NewJournalEntryDescription(str string) (JournalEntryDescription, error) {
	jed := JournalEntryDescription(str)
	if err := jed.Validate(); err != nil {
		return "", err
	}
	return jed, nil
}

// Validate whether the journal entry description is valid. Valid is if:
//   - non empty
//   - valid utf8 encoding
//   - does not have leading or trailing whitespaces
//   - between configured length bounds
//   - contains no control characters
func (jed JournalEntryDescription) Validate() error {
	if err := validateText(string(jed), journalEntryDescriptionMinLen, journalEntryDescriptionMaxLen); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidJournalEntryDescription, err)
	}
	return nil
}
