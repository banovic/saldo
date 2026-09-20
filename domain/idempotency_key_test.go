package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestIdempotencyKeyValidate(t *testing.T) {
	testCases := []struct {
		name           string
		idempotencyKey IdempotencyKey
		wantErr        error
	}{
		{"empty", "", ErrInvalidIdempotencyKey},
		{"typical key", "order-4711", nil},
		{"uuid shaped", "01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ea0", nil},
		{"underscores", "invoice_2026_09", nil},
		{"all allowed classes", "aZ0_-", nil},

		{"2 chars is too short", "ab", ErrInvalidIdempotencyKey},
		{"3 chars is the minimum", "abc", nil},
		{"100 chars is the maximum", IdempotencyKey(strings.Repeat("a", 100)), nil},
		{"101 chars is too long", IdempotencyKey(strings.Repeat("a", 101)), ErrInvalidIdempotencyKey},

		// Length is counted in bytes, but only ASCII passes the character check,
		// so multi-byte input is rejected either way.
		{"multi-byte chars", "Šđč", ErrInvalidIdempotencyKey},
		{"100 multi-byte chars", IdempotencyKey(strings.Repeat("ž", 100)), ErrInvalidIdempotencyKey},

		{"space", "key 1", ErrInvalidIdempotencyKey},
		{"leading whitespace", " key", ErrInvalidIdempotencyKey},
		{"trailing whitespace", "key ", ErrInvalidIdempotencyKey},
		{"dot", "key.1", ErrInvalidIdempotencyKey},
		{"slash", "key/1", ErrInvalidIdempotencyKey},
		{"colon", "key:1", ErrInvalidIdempotencyKey},
		{"control character", "key\x00", ErrInvalidIdempotencyKey},
		{"newline", "key\n1", ErrInvalidIdempotencyKey},
		{"invalid utf8", IdempotencyKey([]byte{0xff, 0xfe, 0xfd}), ErrInvalidIdempotencyKey},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.idempotencyKey.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.idempotencyKey, err, tc.wantErr)
			}
		})
	}
}

func TestNewIdempotencyKey(t *testing.T) {
	testCases := []struct {
		name    string
		str     string
		want    IdempotencyKey
		wantErr error
	}{
		{"valid", "order-4711", "order-4711", nil},
		{"invalid returns zero value", "key 1", "", ErrInvalidIdempotencyKey},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewIdempotencyKey(tc.str)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewIdempotencyKey(%q) error = %v, want %v", tc.str, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("NewIdempotencyKey(%q) = %q, want %q", tc.str, got, tc.want)
			}
		})
	}
}
