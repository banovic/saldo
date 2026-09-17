package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestLedgerNameValidate(t *testing.T) {
	testCases := []struct {
		name       string
		ledgerName LedgerName
		wantErr    error
	}{
		{"empty", "", ErrInvalidLedgerName},
		{"2 chars is too short", "ab", ErrInvalidLedgerName},
		{"3 chars is the minimum", "abc", nil},
		{"100 chars is the maximum", LedgerName(strings.Repeat("a", 100)), nil},
		{"101 chars is too long", LedgerName(strings.Repeat("a", 101)), ErrInvalidLedgerName},

		// Length is counted in characters, not bytes.
		{"2 multi-byte chars is too short, although 4 bytes", "Šđ", ErrInvalidLedgerName},
		{"3 multi-byte chars", "Šđč", nil},
		{"100 multi-byte chars, although 200 bytes", LedgerName(strings.Repeat("ž", 100)), nil},
		{"101 multi-byte chars", LedgerName(strings.Repeat("ž", 101)), ErrInvalidLedgerName},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.ledgerName.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.ledgerName, err, tc.wantErr)
			}
		})
	}
}
