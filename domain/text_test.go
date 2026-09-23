package domain

import (
	"strings"
	"testing"
)

func TestValidateText(t *testing.T) {
	testCases := []struct {
		name    string
		s       string
		min     int
		max     int
		wantErr bool
	}{
		{"plain text", "Office rent", 3, 100, false},
		{"inner spaces are fine", "a b  c", 3, 100, false},
		{"empty", "", 0, 100, true},
		{"invalid utf8", string([]byte{0xff, 0xfe, 'a'}), 1, 100, true},

		{"leading space", " abc", 1, 100, true},
		{"trailing space", "abc ", 1, 100, true},
		{"leading tab", "\tabc", 1, 100, true},
		{"trailing newline", "abc\n", 1, 100, true},

		{"below min", "ab", 3, 100, true},
		{"at min", "abc", 3, 100, false},
		{"at max", strings.Repeat("a", 100), 3, 100, false},
		{"above max", strings.Repeat("a", 101), 3, 100, true},

		// Length is counted in characters, not bytes.
		{"multi-byte chars below min", "Šđ", 3, 100, true},
		{"multi-byte chars at min", "Šđč", 3, 100, false},
		{"multi-byte chars at max, although 200 bytes", strings.Repeat("ž", 100), 3, 100, false},
		{"multi-byte chars above max", strings.Repeat("ž", 101), 3, 100, true},

		{"inner newline", "a\nb", 1, 100, true},
		{"inner tab", "a\tb", 1, 100, true},
		{"inner null", "a\x00b", 1, 100, true},
		{"inner delete", "a\u007fb", 1, 100, true},
		{"emoji is not a control char", "a🙂b", 1, 100, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateText(tc.s, tc.min, tc.max)
			if got := err != nil; got != tc.wantErr {
				t.Errorf("validateText(%q, %d, %d) = %v, want error: %t", tc.s, tc.min, tc.max, err, tc.wantErr)
			}
		})
	}
}
