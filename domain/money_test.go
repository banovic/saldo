package domain

import (
	"errors"
	"maps"
	"math"
	"slices"
	"testing"
)

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

func TestSum(t *testing.T) {
	testCases := []struct {
		name     string
		currency Currency
		ms       []Money
		want     Money
		wantErr  error
	}{
		{"empty currency", "", nil, Money{}, ErrInvalidCurrency},
		{"unknown currency", "XXX", nil, Money{}, ErrInvalidCurrency},
		{"unknown currency is checked before the monies", "XXX", []Money{{1, USD}}, Money{}, ErrInvalidCurrency},

		{"nil slice is zero", USD, nil, Money{0, USD}, nil},
		{"empty slice is zero", USD, []Money{}, Money{0, USD}, nil},

		{"single money", USD, []Money{{500, USD}}, Money{500, USD}, nil},
		{"two monies", USD, []Money{{500, USD}, {250, USD}}, Money{750, USD}, nil},
		{"negatives", USD, []Money{{-500, USD}, {-250, USD}}, Money{-750, USD}, nil},
		{"mixed signs", USD, []Money{{500, USD}, {-250, USD}}, Money{250, USD}, nil},
		{"balanced entry sums to zero", USD, []Money{{500, USD}, {-200, USD}, {-300, USD}}, Money{0, USD}, nil},
		{"zero exponent currency", JPY, []Money{{5, JPY}, {7, JPY}}, Money{12, JPY}, nil},

		{"mismatched currency", USD, []Money{{500, USD}, {250, EUR}}, Money{}, ErrCurrencyMismatch},
		{"mismatch on first money", USD, []Money{{250, EUR}, {500, USD}}, Money{}, ErrCurrencyMismatch},
		{"invalid currency in a money is a mismatch", USD, []Money{{500, USD}, {250, "XXX"}}, Money{}, ErrCurrencyMismatch},
		{"empty currency in a money is a mismatch", USD, []Money{{500, USD}, {250, ""}}, Money{}, ErrCurrencyMismatch},

		{"sum exactly MaxInt64", USD, []Money{{math.MaxInt64 - 1, USD}, {1, USD}}, Money{math.MaxInt64, USD}, nil},
		{"sum exactly MinInt64", USD, []Money{{math.MinInt64 + 1, USD}, {-1, USD}}, Money{math.MinInt64, USD}, nil},
		{"one past MaxInt64 overflows", USD, []Money{{math.MaxInt64, USD}, {1, USD}}, Money{}, ErrOverflow},
		{"one past MinInt64 overflows", USD, []Money{{math.MinInt64, USD}, {-1, USD}}, Money{}, ErrOverflow},
		{"max plus max overflows", USD, []Money{{math.MaxInt64, USD}, {math.MaxInt64, USD}}, Money{}, ErrOverflow},
		{"min plus min overflows", USD, []Money{{math.MinInt64, USD}, {math.MinInt64, USD}}, Money{}, ErrOverflow},

		// A running int64 total would overflow on the second money, but the
		// total is representable, so Sum must succeed regardless of order.
		{"partial sum overflows, total fits", USD,
			[]Money{{math.MaxInt64, USD}, {math.MaxInt64, USD}, {math.MinInt64, USD}},
			Money{math.MaxInt64 - 1, USD}, nil},
		{"same monies reordered", USD,
			[]Money{{math.MaxInt64, USD}, {math.MinInt64, USD}, {math.MaxInt64, USD}},
			Money{math.MaxInt64 - 1, USD}, nil},
		{"partial sum underflows, total fits", USD,
			[]Money{{math.MinInt64, USD}, {math.MinInt64, USD}, {math.MaxInt64, USD}, {math.MaxInt64, USD}},
			Money{-2, USD}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SumMoney(tc.currency, tc.ms)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Sum(%q, %v) error = %v, want %v", tc.currency, tc.ms, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Sum(%q, %v) = %v, want %v", tc.currency, tc.ms, got, tc.want)
			}
		})
	}
}

func TestSumDoesNotModifyInput(t *testing.T) {
	ms := []Money{{500, USD}, {-200, USD}}
	want := []Money{{500, USD}, {-200, USD}}
	if _, err := SumMoney(USD, ms); err != nil {
		t.Fatalf("Sum(%q, %v) error = %v, want nil", USD, ms, err)
	}
	for i := range ms {
		if ms[i] != want[i] {
			t.Errorf("Sum modified ms[%d] = %v, want %v", i, ms[i], want[i])
		}
	}
}

func TestMoneyIsZeroAmount(t *testing.T) {
	testCases := []struct {
		name string
		a    Money
		want bool
	}{
		{"zero value money", Money{}, true},
		{"zero amount", Money{0, USD}, true},
		{"currency is ignored", Money{0, "XXX"}, true},
		{"zero in a currency without minor units", Money{0, JPY}, true},
		{"positive", Money{1, USD}, false},
		{"negative", Money{-1, USD}, false},
		{"MaxInt64", Money{math.MaxInt64, USD}, false},
		{"MinInt64", Money{math.MinInt64, USD}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.a.IsZeroAmount(); got != tc.want {
				t.Errorf("%v.IsZeroAmount() = %t, want %t", tc.a, got, tc.want)
			}
		})
	}
}

// currencyInfo is hand-maintained data. These tests walk the table itself
// rather than a list of cases, so a currency added later is checked too.

func TestCurrencyInfoTableShape(t *testing.T) {
	for _, c := range slices.Sorted(maps.Keys(currencyInfo)) {
		ci := currencyInfo[c]
		t.Run(string(c), func(t *testing.T) {
			if len(c) != 3 {
				t.Errorf("code %q has %d characters, want 3", c, len(c))
			}
			for _, r := range c {
				if r < 'A' || r > 'Z' {
					t.Errorf("code %q contains %q, want uppercase A-Z", c, r)
				}
			}
			if ci.Name == "" {
				t.Errorf("%q has no Name", c)
			}
			if len(ci.Num) != 3 {
				t.Errorf("%q has Num %q with %d characters, want 3", c, ci.Num, len(ci.Num))
			}
			for _, r := range ci.Num {
				if r < '0' || r > '9' {
					t.Errorf("%q has Num %q containing %q, want digits", c, ci.Num, r)
				}
			}
			// Guards the unreachable branch in MinorUnitsPerUnit: it is only
			// unreachable as long as every exponent in the table is supported.
			if _, ok := ci.MinorUnitsPerUnit(); !ok {
				t.Errorf("%q has Exponent %d, which MinorUnitsPerUnit does not support", c, ci.Exponent)
			}
		})
	}
}

func TestCurrencyInfoNumIsUnique(t *testing.T) {
	seen := make(map[string]Currency)
	for _, c := range slices.Sorted(maps.Keys(currencyInfo)) {
		num := currencyInfo[c].Num
		if prev, ok := seen[num]; ok {
			t.Errorf("Num %q is used by both %q and %q", num, prev, c)
			continue
		}
		seen[num] = c
	}
}

// A Currency constant with no row in currencyInfo would be silently invalid:
// EUR.IsValid() would report false. The length check catches the reverse,
// a row added to the table with no constant declared for it.
func TestCurrencyConstantsHaveTableRows(t *testing.T) {
	constants := []Currency{EUR, RSD, USD, JPY, TND}
	for _, c := range constants {
		if !c.IsValid() {
			t.Errorf("%q is declared as a constant but has no row in currencyInfo", c)
		}
	}
	if len(constants) != len(currencyInfo) {
		t.Errorf("%d Currency constants but %d rows in currencyInfo", len(constants), len(currencyInfo))
	}
}
