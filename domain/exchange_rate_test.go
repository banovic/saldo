package domain

import (
	"testing"
	"time"
)

func TestExchangeRateIsZero(t *testing.T) {
	testCases := []struct {
		name string
		er   ExchangeRate
		want bool
	}{
		{"zero value", ExchangeRate{}, true},
		{"only ID set", ExchangeRate{ExchangeRateID: "r1"}, false},
		{"only Num set", ExchangeRate{Num: 1}, false},
		{"only On set", ExchangeRate{On: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)}, false},
		{"fully set", ExchangeRate{ExchangeRateID: "r1", From: EUR, To: RSD, Num: 1178, Den: 10, Source: NBS, Kind: Spot}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.er.IsZero(); got != tc.want {
				t.Errorf("%+v.IsZero() = %t, want %t", tc.er, got, tc.want)
			}
		})
	}
}
