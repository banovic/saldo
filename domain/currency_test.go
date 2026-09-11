package domain

import (
	"testing"
)

func TestCurrencyIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		currency Currency
		want     bool
	}{
		{"empty currency", "", false},
		{"invalid currency", "X12", false},
		{"invalid currency 2, ISO 4217 codes are uppercase", "usd", false},
		{"valid currency", USD, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.currency.IsValid(); got != tc.want {
				t.Errorf("%q.IsValid() = %t, want %t", tc.currency, got, tc.want)
			}
		})
	}
}

func TestCurrencyInfo(t *testing.T) {
	testCases := []struct {
		name     string
		currency Currency
		want     CurrencyInfo
		wantOk   bool
	}{
		{"empty currency", "", CurrencyInfo{}, false},
		{"invalid currency", "X12", CurrencyInfo{}, false},
		{"valid currency", USD, CurrencyInfo{Name: "US Dollar", Num: "840", Exponent: 2}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.currency.Info()
			if ok != tc.wantOk || got != tc.want {
				t.Errorf("%q.Info() = %v, %t, want %v, %t", tc.currency, got, ok, tc.want, tc.wantOk)
			}
		})
	}
}
