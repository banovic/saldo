package domain

import (
	"errors"
	"math"
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

func TestExchangeRateValidate(t *testing.T) {
	valid := ExchangeRate{
		ExchangeRateID: ExchangeRateID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8eb0")},
		From:           EUR,
		To:             RSD,
		Num:            11780,
		Den:            100,
		On:             Date{2026, time.September, 18},
	}
	with := func(f func(*ExchangeRate)) ExchangeRate {
		er := valid
		f(&er)
		return er
	}

	testCases := []struct {
		name    string
		rate    ExchangeRate
		wantErr error
	}{
		{"valid", valid, nil},
		{"rate of one", with(func(er *ExchangeRate) { er.Num, er.Den = 1, 1 }), nil},

		{"zero exchange rate id", with(func(er *ExchangeRate) { er.ExchangeRateID = ExchangeRateID{} }), ErrInvalidExchangeRate},
		{"empty from currency", with(func(er *ExchangeRate) { er.From = "" }), ErrInvalidCurrency},
		{"unknown from currency", with(func(er *ExchangeRate) { er.From = "XXX" }), ErrInvalidCurrency},
		{"empty to currency", with(func(er *ExchangeRate) { er.To = "" }), ErrInvalidCurrency},
		{"unknown to currency", with(func(er *ExchangeRate) { er.To = "XXX" }), ErrInvalidCurrency},
		{"same currencies", with(func(er *ExchangeRate) { er.To = er.From }), ErrInvalidExchangeRate},
		{"zero numerator", with(func(er *ExchangeRate) { er.Num = 0 }), ErrInvalidExchangeRate},
		{"negative numerator", with(func(er *ExchangeRate) { er.Num = -1 }), ErrInvalidExchangeRate},
		{"zero denominator", with(func(er *ExchangeRate) { er.Den = 0 }), ErrInvalidExchangeRate},
		{"negative denominator", with(func(er *ExchangeRate) { er.Den = -1 }), ErrInvalidExchangeRate},
		{"zero on date", with(func(er *ExchangeRate) { er.On = Date{} }), ErrInvalidDate},
		{"impossible on date", with(func(er *ExchangeRate) { er.On = Date{2026, time.February, 30} }), ErrInvalidDate},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.rate.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%+v.Validate() = %v, want %v", tc.rate, err, tc.wantErr)
			}
			if tc.wantErr != nil && !errors.Is(err, ErrInvalidExchangeRate) {
				t.Errorf("%+v.Validate() = %v, want %v", tc.rate, err, ErrInvalidExchangeRate)
			}
		})
	}
}

func TestExchangeRateConvert(t *testing.T) {
	rateID := ExchangeRateID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8eb0")}
	rate := func(from, to Currency, num, den int64) ExchangeRate {
		return ExchangeRate{
			ExchangeRateID: rateID,
			From:           from,
			To:             to,
			Num:            num,
			Den:            den,
			On:             Date{2026, time.September, 18},
		}
	}

	// 1 EUR = 117.8 RSD, and its inverse.
	eurRSD := rate(EUR, RSD, 1178, 10)
	rsdEUR := rate(RSD, EUR, 10, 1178)
	// Rates chosen so that conversion lands exactly on a tie, or just off one.
	half := rate(USD, JPY, 1, 2)
	third := rate(USD, JPY, 1, 3)
	identity := rate(EUR, USD, 1, 1)

	testCases := []struct {
		name    string
		er      ExchangeRate
		m       Money
		want    Money
		wantErr error
	}{
		{"exact conversion", eurRSD, Money{100, EUR}, Money{11780, RSD}, nil},
		{"zero amount", eurRSD, Money{0, EUR}, Money{0, RSD}, nil},
		{"negative amount", eurRSD, Money{-100, EUR}, Money{-11780, RSD}, nil},
		{"exact conversion, larger amount", eurRSD, Money{12345, EUR}, Money{1454241, RSD}, nil},
		{"remainder rounds up", eurRSD, Money{1, EUR}, Money{118, RSD}, nil},
		{"remainder rounds up, negative", eurRSD, Money{-1, EUR}, Money{-118, RSD}, nil},
		{"remainder rounds up, two minor units", eurRSD, Money{2, EUR}, Money{236, RSD}, nil},
		{"remainder divides evenly", eurRSD, Money{5, EUR}, Money{589, RSD}, nil},

		// Halves go away from zero, matching divAndRoundHalfUp.
		{"tie rounds away from zero", half, Money{1, USD}, Money{1, JPY}, nil},
		{"tie rounds away from zero, negative", half, Money{-1, USD}, Money{-1, JPY}, nil},
		{"tie above one", half, Money{3, USD}, Money{2, JPY}, nil},
		{"tie above one, negative", half, Money{-3, USD}, Money{-2, JPY}, nil},
		{"no tie, exact", half, Money{4, USD}, Money{2, JPY}, nil},

		{"rounds down to zero", third, Money{1, USD}, Money{0, JPY}, nil},
		{"rounds down to zero, negative", third, Money{-1, USD}, Money{0, JPY}, nil},
		{"rounds up from a third", third, Money{2, USD}, Money{1, JPY}, nil},
		{"rounds up from a third, negative", third, Money{-2, USD}, Money{-1, JPY}, nil},

		// Just either side of a tie, with a denominator that is not a power of ten.
		{"just above a tie", rsdEUR, Money{59, RSD}, Money{1, EUR}, nil},
		{"just below a tie", rsdEUR, Money{58, RSD}, Money{0, EUR}, nil},
		{"inverse rate is exact", rsdEUR, Money{11780, RSD}, Money{100, EUR}, nil},

		// Num == Den == 1 keeps the amount and only changes the currency.
		{"rate of one", identity, Money{12345, EUR}, Money{12345, USD}, nil},
		{"rate of one, max int64", identity, Money{math.MaxInt64, EUR}, Money{math.MaxInt64, USD}, nil},
		{"rate of one, min int64", identity, Money{math.MinInt64, EUR}, Money{math.MinInt64, USD}, nil},

		{"money is not in the rate's from currency", eurRSD, Money{100, USD}, Money{}, ErrCurrencyMismatch},
		{"money has no currency", eurRSD, Money{100, ""}, Money{}, ErrCurrencyMismatch},
		{"rate is used in the wrong direction", eurRSD, Money{11780, RSD}, Money{}, ErrCurrencyMismatch},
		{"zero amount in the wrong currency still fails", eurRSD, Money{0, USD}, Money{}, ErrCurrencyMismatch},

		// Num * MinorUnits is computed before the division, so it can overflow
		// even when the converted amount would fit.
		{"numerator times amount overflows", eurRSD, Money{math.MaxInt64, EUR}, Money{}, ErrOverflow},
		{"numerator times amount underflows", eurRSD, Money{math.MinInt64, EUR}, Money{}, ErrOverflow},
		{"overflow wraps the conversion error", eurRSD, Money{math.MaxInt64, EUR}, Money{}, ErrInvalidExchangeRateConvert},

		// Convert assumes a valid rate, but a zero denominator is reported
		// rather than panicking.
		{"zero denominator", rate(EUR, RSD, 1178, 0), Money{100, EUR}, Money{}, ErrInvalidExchangeRateConvert},
		{"negative denominator", rate(EUR, RSD, 1178, -10), Money{100, EUR}, Money{}, ErrInvalidExchangeRateConvert},
		{"zero numerator converts to zero", rate(EUR, RSD, 0, 10), Money{100, EUR}, Money{0, RSD}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.er.Convert(tc.m)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%v.Convert(%v) error = %v, want %v", tc.er, tc.m, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("%v.Convert(%v) = %v, want %v", tc.er, tc.m, got, tc.want)
			}
		})
	}
}
