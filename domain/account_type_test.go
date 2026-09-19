package domain

import (
	"errors"
	"testing"
)

func TestAccountTypeValidate(t *testing.T) {
	testCases := []struct {
		name        string
		accountType AccountType
		wantErr     error
	}{
		{"empty account type", "", ErrInvalidAccountType},
		{"unknown account type", "income", ErrInvalidAccountType},
		{"account types are lowercase", "Asset", ErrInvalidAccountType},
		{"account types are lowercase 2", "ASSET", ErrInvalidAccountType},
		{"whitespace is not trimmed", " asset", ErrInvalidAccountType},
		{"valid account type, asset", Asset, nil},
		{"valid account type, liability", Liability, nil},
		{"valid account type, equity", Equity, nil},
		{"valid account type, revenue", Revenue, nil},
		{"valid account type, expense", Expense, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.accountType.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.accountType, err, tc.wantErr)
			}
		})
	}
}
