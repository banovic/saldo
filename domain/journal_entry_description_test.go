package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestJournalEntryDescriptionValidate(t *testing.T) {
	testCases := []struct {
		name                    string
		journalEntryDescription JournalEntryDescription
		wantErr                 error
	}{
		{"empty", "", ErrInvalidJournalEntryDescription},
		{"2 chars is too short", "ab", ErrInvalidJournalEntryDescription},
		{"3 chars is the minimum", "abc", nil},
		{"100 chars is the maximum", JournalEntryDescription(strings.Repeat("a", 100)), nil},
		{"101 chars is too long", JournalEntryDescription(strings.Repeat("a", 101)), ErrInvalidJournalEntryDescription},

		// Length is counted in characters, not bytes.
		{"2 multi-byte chars is too short, although 4 bytes", "Šđ", ErrInvalidJournalEntryDescription},
		{"3 multi-byte chars", "Šđč", nil},
		{"100 multi-byte chars, although 200 bytes", JournalEntryDescription(strings.Repeat("ž", 100)), nil},
		{"101 multi-byte chars", JournalEntryDescription(strings.Repeat("ž", 101)), ErrInvalidJournalEntryDescription},

		{"inner whitespace", "Office rent", nil},
		{"leading space", " abc", ErrInvalidJournalEntryDescription},
		{"trailing newline", "abc\n", ErrInvalidJournalEntryDescription},
		{"control char", "ab\x00c", ErrInvalidJournalEntryDescription},
		{"invalid utf8", "ab\xffc", ErrInvalidJournalEntryDescription},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.journalEntryDescription.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.journalEntryDescription, err, tc.wantErr)
			}
		})
	}
}
