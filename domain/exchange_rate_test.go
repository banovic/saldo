package domain

import (
	"testing"
	"time"
	"uuid"
)

func TestExchangeRateIDIsZero(t *testing.T) {
	testCases := []struct {
		name string
		id   ExchangeRateID
		want bool
	}{
		{"zero value", ExchangeRateID{}, true},
		{"nil uuid", ExchangeRateID{uuid.Nil()}, true},
		{"non-nil uuid", ExchangeRateID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8eb0")}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.IsZero(); got != tc.want {
				t.Errorf("%v.IsZero() = %t, want %t", tc.id, got, tc.want)
			}
		})
	}
}

func TestExchangeRateIsZero(t *testing.T) {
	rateID := ExchangeRateID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8eb0")}

	testCases := []struct {
		name string
		er   ExchangeRate
		want bool
	}{
		{"zero value", ExchangeRate{}, true},
		{"only ID set", ExchangeRate{ExchangeRateID: rateID}, false},
		{"only Num set", ExchangeRate{Num: 1}, false},
		{"only On set", ExchangeRate{On: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)}, false},
		{"fully set", ExchangeRate{ExchangeRateID: rateID, From: EUR, To: RSD, Num: 1178, Den: 10, Source: NBS, Kind: Spot}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.er.IsZero(); got != tc.want {
				t.Errorf("%+v.IsZero() = %t, want %t", tc.er, got, tc.want)
			}
		})
	}
}
