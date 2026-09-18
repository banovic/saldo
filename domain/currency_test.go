package domain

import (
	"errors"
	"testing"
)

func TestCurrencyValidate(t *testing.T) {
	testCases := []struct {
		name     string
		currency Currency
		wantErr  error
	}{
		{"empty currency", "", ErrInvalidCurrency},
		{"invalid currency", "X12", ErrInvalidCurrency},
		{"invalid currency 2, ISO 4217 codes are uppercase", "usd", ErrInvalidCurrency},
		{"valid currency", USD, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.currency.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.currency, err, tc.wantErr)
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
		{"valid currency", USD, CurrencyInfo{Name: "US Dollar", Num: "840", Exponent: 2, MinorUnitsInUnit: 100}, true},
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
