package domain

import (
	"errors"
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
