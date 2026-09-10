package saldo

import "testing"

func TestAccountTypeIsValid(t *testing.T) {
	testCases := []struct {
		name        string
		accountType AccountType
		want        bool
	}{
		{"empty account type", "", false},
		{"unknown account type", "income", false},
		{"account types are lowercase", "Asset", false},
		{"account types are lowercase 2", "ASSET", false},
		{"whitespace is not trimmed", " asset", false},
		{"valid account type, asset", Asset, true},
		{"valid account type, liability", Liability, true},
		{"valid account type, equity", Equity, true},
		{"valid account type, revenue", Revenue, true},
		{"valid account type, expense", Expense, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.accountType.IsValid(); got != tc.want {
				t.Errorf("%q.IsValid() = %t, want %t", tc.accountType, got, tc.want)
			}
		})
	}
}
