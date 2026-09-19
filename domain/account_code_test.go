package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestAccountCodeValidate(t *testing.T) {
	testCases := []struct {
		name        string
		accountCode AccountCode
		wantErr     error
	}{
		{"empty", "", ErrInvalidAccountCode},
		{"typical code", "241", nil},
		{"alphanumeric", "A-100", nil},

		{"2 chars is too short", "24", ErrInvalidAccountCode},
		{"3 chars is the minimum", "241", nil},
		{"100 chars is the maximum", AccountCode(strings.Repeat("1", 100)), nil},
		{"101 chars is too long", AccountCode(strings.Repeat("1", 101)), ErrInvalidAccountCode},

		// Length is counted in characters, not bytes.
		{"2 multi-byte chars is too short, although 4 bytes", "Šđ", ErrInvalidAccountCode},
		{"3 multi-byte chars", "Šđč", nil},
		{"100 multi-byte chars, although 200 bytes", AccountCode(strings.Repeat("ž", 100)), nil},
		{"101 multi-byte chars", AccountCode(strings.Repeat("ž", 101)), ErrInvalidAccountCode},

		{"invalid utf8", AccountCode([]byte{0xff, 0xfe, 0xfd}), ErrInvalidAccountCode},
		{"leading whitespace", " 241", ErrInvalidAccountCode},
		{"trailing whitespace", "241 ", ErrInvalidAccountCode},
		{"inner whitespace is allowed", "241 A", nil},
		{"control character", "24\x00", ErrInvalidAccountCode},
		{"newline", "2\n41", ErrInvalidAccountCode},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.accountCode.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.accountCode, err, tc.wantErr)
			}
		})
	}
}
