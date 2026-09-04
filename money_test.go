package saldo

import (
	"errors"
	"math"
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

func TestNewMoney(t *testing.T) {
	testCases := []struct {
		name     string
		units    int64
		currency Currency
		want     Money
		wantErr  error
	}{
		{"empty currency", 0, "", Money{}, ErrInvalidCurrency},
		{"invalid currency", 0, "X12", Money{}, ErrInvalidCurrency},
		{"invalid currency 2", 10, "X12", Money{}, ErrInvalidCurrency},
		{"valid currency", 10, RSD, Money{MinorUnits: 10, Currency: RSD}, nil},
		{"valid currency 2", -10, EUR, Money{MinorUnits: -10, Currency: EUR}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewMoney(tc.units, tc.currency)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewMoney(%d, %q) error = %v, want %v", tc.units, tc.currency, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("NewMoney(%d, %q) = %v want %v", tc.units, tc.currency, got, tc.want)
			}
		})
	}
}

func TestZeroMoney(t *testing.T) {
	testCases := []struct {
		name     string
		currency Currency
		want     Money
		wantErr  error
	}{
		{"empty currency", "", Money{}, ErrInvalidCurrency},
		{"invalid currency", "X12", Money{}, ErrInvalidCurrency},
		{"invalid currency 2", "usd", Money{}, ErrInvalidCurrency},
		{"valid currency", RSD, Money{MinorUnits: 0, Currency: RSD}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ZeroMoney(tc.currency)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ZeroMoney(%q) error = %v, want %v", tc.currency, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("ZeroMoney(%q) = %v want %v", tc.currency, got, tc.want)
			}
		})
	}
}

func TestMoneyAdd(t *testing.T) {
	testCases := []struct {
		name    string
		a       Money
		b       Money
		want    Money
		wantErr error
	}{
		{"zero plus zero is error", Money{}, Money{}, Money{}, ErrInvalidMoney},
		{"zero plus valid is error", Money{}, Money{MinorUnits: 1, Currency: USD}, Money{}, ErrInvalidMoney},
		{"valid plus zero is error", Money{MinorUnits: 1, Currency: USD}, Money{}, Money{}, ErrInvalidMoney},
		{"currency mismatch", Money{MinorUnits: -1, Currency: USD}, Money{MinorUnits: 1, Currency: EUR}, Money{}, ErrCurrencyMismatch},
		{"sum exactly MaxInt64", Money{1, USD}, Money{math.MaxInt64 - 1, USD}, Money{math.MaxInt64, USD}, nil},
		{"sum exactly MinInt64", Money{-1, USD}, Money{math.MinInt64 + 1, USD}, Money{math.MinInt64, USD}, nil},
		{"max plus zero", Money{math.MaxInt64, USD}, Money{0, USD}, Money{math.MaxInt64, USD}, nil},
		{"zero plus max", Money{0, USD}, Money{math.MaxInt64, USD}, Money{math.MaxInt64, USD}, nil},
		{"max plus max", Money{math.MaxInt64, USD}, Money{math.MaxInt64, USD}, Money{}, ErrOverflow},
		{"min plus min", Money{math.MinInt64, USD}, Money{math.MinInt64, USD}, Money{}, ErrOverflow},
		{"small positive overflows", Money{2, USD}, Money{math.MaxInt64 - 1, USD}, Money{}, ErrOverflow},
		{"small negative underflows", Money{-2, USD}, Money{math.MinInt64 + 1, USD}, Money{}, ErrOverflow},
		{"ok", Money{MinorUnits: -1, Currency: USD}, Money{MinorUnits: 1, Currency: USD}, Money{MinorUnits: 0, Currency: USD}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%v.Add(%v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("%v.Add(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMoneySub(t *testing.T) {
	testCases := []struct {
		name    string
		a       Money
		b       Money
		want    Money
		wantErr error
	}{
		{"zero minus zero is error", Money{}, Money{}, Money{}, ErrInvalidMoney},
		{"zero minus valid is error", Money{}, Money{MinorUnits: 1, Currency: USD}, Money{}, ErrInvalidMoney},
		{"valid minus zero is error", Money{MinorUnits: 1, Currency: USD}, Money{}, Money{}, ErrInvalidMoney},
		{"currency mismatch", Money{MinorUnits: -1, Currency: USD}, Money{MinorUnits: 1, Currency: EUR}, Money{}, ErrCurrencyMismatch},
		{"diff exactly MaxInt64", Money{math.MaxInt64 - 1, USD}, Money{-1, USD}, Money{math.MaxInt64, USD}, nil},
		{"diff exactly MinInt64", Money{math.MinInt64 + 1, USD}, Money{1, USD}, Money{math.MinInt64, USD}, nil},
		{"just over MaxInt64", Money{math.MaxInt64, USD}, Money{-1, USD}, Money{}, ErrOverflow},
		{"just under MinInt64", Money{math.MinInt64, USD}, Money{1, USD}, Money{}, ErrOverflow},
		{"zero minus min overflows", Money{0, USD}, Money{math.MinInt64, USD}, Money{}, ErrOverflow},
		{"ok", Money{-1, USD}, Money{1, USD}, Money{-2, USD}, nil},
		{"ok", Money{MinorUnits: -1, Currency: USD}, Money{MinorUnits: 1, Currency: USD}, Money{MinorUnits: -2, Currency: USD}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.a.Sub(tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%v.Sub(%v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("%v.Sub(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMoneyMul(t *testing.T) {
	testCases := []struct {
		name    string
		a       Money
		k       int64
		want    Money
		wantErr error
	}{
		{"invalid money", Money{5, "XXX"}, 3, Money{}, ErrInvalidMoney},
		{"times zero", Money{500, USD}, 0, Money{0, USD}, nil},
		{"zero times k", Money{0, USD}, 7, Money{0, USD}, nil},
		{"identity", Money{500, USD}, 1, Money{500, USD}, nil},
		{"negate", Money{500, USD}, -1, Money{-500, USD}, nil},
		{"normal", Money{250, USD}, 4, Money{1000, USD}, nil},
		{"two negatives", Money{-250, USD}, -4, Money{1000, USD}, nil},
		{"MinInt64 times -1 overflows", Money{math.MinInt64, USD}, -1, Money{}, ErrOverflow},
		{"MinInt64 times 1 fits", Money{math.MinInt64, USD}, 1, Money{math.MinInt64, USD}, nil},
		{"MaxInt64 times -1 fits", Money{math.MaxInt64, USD}, -1, Money{math.MinInt64 + 1, USD}, nil},
		{"MaxInt64 times 2 overflows", Money{math.MaxInt64, USD}, 2, Money{}, ErrOverflow},
		{"minus one times MinInt64 overflows", Money{-1, USD}, math.MinInt64, Money{}, ErrOverflow},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.a.Mul(tc.k)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%v.Mul(%d) error = %v, want %v", tc.a, tc.k, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("%v.Mul(%d) = %v, want %v", tc.a, tc.k, got, tc.want)
			}
		})
	}
}

func TestMoneyString(t *testing.T) {
	testCases := []struct {
		name string
		a    Money
		want string
	}{
		{"whole units", Money{1050, USD}, "10.50 USD"},
		{"pads fraction", Money{1005, USD}, "10.05 USD"},
		{"zero fraction", Money{1000, USD}, "10.00 USD"},
		{"under one unit", Money{5, USD}, "0.05 USD"},
		{"zero", Money{0, USD}, "0.00 USD"},
		{"negative", Money{-1050, USD}, "-10.50 USD"},
		{"negative pads", Money{-1005, USD}, "-10.05 USD"},
		{"negative under one", Money{-5, USD}, "-0.05 USD"},
		{"MinInt64", Money{math.MinInt64, USD}, "-92233720368547758.08 USD"},
		{"MaxInt64", Money{math.MaxInt64, USD}, "92233720368547758.07 USD"},
		{"zero money", Money{}, `0 minor units of ""`},

		// TND: exponent 3
		{"TND whole units", Money{1050, TND}, "1.050 TND"},
		{"TND exact unit", Money{1000, TND}, "1.000 TND"},
		{"TND under one unit", Money{5, TND}, "0.005 TND"},
		{"TND zero", Money{0, TND}, "0.000 TND"},
		{"TND negative", Money{-1050, TND}, "-1.050 TND"},
		{"TND negative under one", Money{-5, TND}, "-0.005 TND"},
		{"TND MinInt64", Money{math.MinInt64, TND}, "-9223372036854775.808 TND"},

		// JPY: exponent 0, no minor units
		{"JPY whole units", Money{1050, JPY}, "1050 JPY"},
		{"JPY zero", Money{0, JPY}, "0 JPY"},
		{"JPY negative", Money{-1050, JPY}, "-1050 JPY"},
		{"JPY MinInt64", Money{math.MinInt64, JPY}, "-9223372036854775808 JPY"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.String()
			if got != tc.want {
				t.Errorf("%v.String() = %v, want %q", tc.a, got, tc.want)
			}
		})
	}
}
