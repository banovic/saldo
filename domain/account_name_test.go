package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestAccountNameValidate(t *testing.T) {
	testCases := []struct {
		name        string
		accountName AccountName
		wantErr     error
	}{
		{"empty", "", ErrInvalidAccountName},
		{"2 chars is too short", "ab", ErrInvalidAccountName},
		{"3 chars is the minimum", "abc", nil},
		{"100 chars is the maximum", AccountName(strings.Repeat("a", 100)), nil},
		{"101 chars is too long", AccountName(strings.Repeat("a", 101)), ErrInvalidAccountName},

		// Length is counted in characters, not bytes.
		{"2 multi-byte chars is too short, although 4 bytes", "Šđ", ErrInvalidAccountName},
		{"3 multi-byte chars", "Šđč", nil},
		{"100 multi-byte chars, although 200 bytes", AccountName(strings.Repeat("ž", 100)), nil},
		{"101 multi-byte chars", AccountName(strings.Repeat("ž", 101)), ErrInvalidAccountName},

		{"inner whitespace is allowed", "Tekući račun", nil},
		{"leading space", " Cash", ErrInvalidAccountName},
		{"trailing space", "Cash ", ErrInvalidAccountName},
		{"trailing newline", "Cash\n", ErrInvalidAccountName},
		{"only whitespace", "   ", ErrInvalidAccountName},
		{"inner newline", "Ca\nsh", ErrInvalidAccountName},
		{"inner tab", "Ca\tsh", ErrInvalidAccountName},
		{"NUL", "Ca\x00sh", ErrInvalidAccountName},
		{"C1 control", "Cash", ErrInvalidAccountName},
		{"invalid utf8", "Ca\xffsh", ErrInvalidAccountName},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.accountName.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.accountName, err, tc.wantErr)
			}
		})
	}
}

func TestNewAccountName(t *testing.T) {
	testCases := []struct {
		name    string
		s       string
		want    AccountName
		wantErr error
	}{
		{"valid", "Cash", "Cash", nil},
		{"invalid returns zero value", "ab", "", ErrInvalidAccountName},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewAccountName(tc.s)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewAccountName(%q) error = %v, want %v", tc.s, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("NewAccountName(%q) = %q, want %q", tc.s, got, tc.want)
			}
		})
	}
}
