package domain

import (
	"errors"
	"testing"
)

func TestExchangeRateSourceValidate(t *testing.T) {
	testCases := []struct {
		name               string
		exchangeRateSource ExchangeRateSource
		wantErr            error
	}{
		{"empty source", "", ErrInvalidExchangeRateSource},
		{"unknown source", "ecb", ErrInvalidExchangeRateSource},
		{"sources are lowercase", "NBS", ErrInvalidExchangeRateSource},
		{"sources are lowercase 2", "Nbs", ErrInvalidExchangeRateSource},
		{"whitespace is not trimmed", " nbs", ErrInvalidExchangeRateSource},
		{"valid source, nbs", NBS, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.exchangeRateSource.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("%q.Validate() = %v, want %v", tc.exchangeRateSource, err, tc.wantErr)
			}
		})
	}
}
