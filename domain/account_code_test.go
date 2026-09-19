package domain

import (
	"errors"
	"testing"
)

func TestAccountCodeValidate(t *testing.T) {
	testCases := []struct {
		name        string
		accountCode AccountCode
		wantErr     error
	}{
		{"empty", "", ErrInvalidAccountCode},
		{"single digit", "2", nil},
		{"typical code", "241", nil},
		{"alphanumeric", "A-100", nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.accountCode.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.accountCode, err, tc.wantErr)
			}
		})
	}
}
