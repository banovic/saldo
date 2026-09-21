package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
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
	str := string(jed)
	if str == "" {
		return fmt.Errorf("%w: empty string", ErrInvalidJournalEntryDescription)
	}
	if !utf8.ValidString(str) {
		return fmt.Errorf("%w: invalid utf8", ErrInvalidJournalEntryDescription)
	}
	if str != strings.TrimSpace(str) {
		return fmt.Errorf("%w: leading or trailing whitespace", ErrInvalidJournalEntryDescription)
	}
	rc := utf8.RuneCountInString(str)
	if rc < journalEntryDescriptionMinLen || rc > journalEntryDescriptionMaxLen {
		return fmt.Errorf("%w: min %d, max %d, got: %d", ErrInvalidJournalEntryDescription, journalEntryDescriptionMinLen, journalEntryDescriptionMaxLen, rc)
	}
	for _, r := range str {
		if unicode.IsControl(r) {
			return fmt.Errorf("%w: invalid char (control): %U", ErrInvalidJournalEntryDescription, r)
		}
	}
	return nil
}
