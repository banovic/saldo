package domain

import "testing"

func TestAccountIDIsValid(t *testing.T) {
	testCases := []struct {
		name      string
		accountID AccountID
		want      bool
	}{
		{"empty account id", "", false},
		{"single character", "a", true},
		{"typical account id", "acc-1", true},
		{"only emptiness is invalid, whitespace is not", " ", true},
		{"whitespace is not trimmed", " acc-1 ", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.accountID.IsValid(); got != tc.want {
				t.Errorf("%q.IsValid() = %t, want %t", tc.accountID, got, tc.want)
			}
		})
	}
}
