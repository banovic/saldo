package saldo

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

var (
	ErrInvalidCurrency  = errors.New("invalid currency")
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrInvalidMoney     = errors.New("invalid money")
	ErrOverflow         = errors.New("overflow")
)

// Currency is ISO 4217 3 letter alphanumeric code.
type Currency string

// CurrencyInfo holds other ISO 4217 information about currency.
type CurrencyInfo struct {
	Name     string
	Num      string
	Exponent uint8
}

const (
	EUR Currency = "EUR"
	RSD Currency = "RSD"
	USD Currency = "USD"
	JPY Currency = "JPY"
	TND Currency = "TND"
)

var currencyInfo = map[Currency]CurrencyInfo{
	EUR: {Name: "Euro", Num: "978", Exponent: 2},
	RSD: {Name: "Serbian dinar", Num: "941", Exponent: 2},
	USD: {Name: "US Dollar", Num: "840", Exponent: 2},
	JPY: {Name: "Japanese yen", Num: "392", Exponent: 0},
	TND: {Name: "Tunisian dinar", Num: "788", Exponent: 3},
}

// Info returns the CurrencyInfo for c and reports whether c is a defined currency.
func (c Currency) Info() (CurrencyInfo, bool) {
	ci, ok := currencyInfo[c]
	return ci, ok
}

// IsValid checks if Currency is valid.
// Currency is valid:
// - must exists as key in currencyInfo map.
func (c Currency) IsValid() bool {
	_, ok := c.Info()
	return ok
}

func (ci CurrencyInfo) MinorUnitsPerUnit() (uint64, bool) {
	switch ci.Exponent {
	case 0:
		return 1, true
	case 1:
		return 10, true
	case 2:
		return 100, true
	case 3:
		return 1000, true
	case 4:
		return 10000, true
	}
	return 0, false
}

// Money represents amount of money in some currency.
// Amount of money is expressed in minor units of that currency, because the focus of this package
// is to track movement of money (ie. not trading, crypto and other specific uses).
type Money struct {
	MinorUnits int64
	Currency   Currency
}

func (m Money) IsValid() bool {
	return m.Currency.IsValid()
}

func (m Money) IsZeroAmount() bool {
	return m.MinorUnits == 0
}

func (m Money) Mul(k int64) (Money, error) {
	if !m.IsValid() {
		return Money{}, fmt.Errorf("%w: %v", ErrInvalidMoney, m)
	}
	if m.MinorUnits == 0 || k == 0 {
		return Money{MinorUnits: 0, Currency: m.Currency}, nil
	}
	// p/k test below cannot detect these values; see 'Go defines MinInt64 / -1 as MinInt64'
	// math.MinInt64 * -1 = math.MinInt64, but this value has overflown already
	// math.MinInt64 / -1 = math.MinInt64, but this value has overflown already
	if m.MinorUnits == math.MinInt64 && k == -1 {
		return Money{}, fmt.Errorf("%w: %v times %d", ErrOverflow, m, k)
	}
	p := m.MinorUnits * k
	if p/k != m.MinorUnits {
		return Money{}, fmt.Errorf("%w: %v times %d", ErrOverflow, m, k)
	}
	return Money{MinorUnits: p, Currency: m.Currency}, nil
}

// Displaying money amounts to end user usually takes into account user's location.
// This method does not do that - it's just for simple display of Money struct.
func (m Money) String() string {
	ci, ok1 := m.Currency.Info()
	d, ok2 := ci.MinorUnitsPerUnit()
	if !ok1 || !ok2 {
		return fmt.Sprintf("%d minor units of %q", m.MinorUnits, m.Currency)
	}
	sign, u := "", uint64(m.MinorUnits)
	if m.MinorUnits < 0 {
		sign, u = "-", -u
	}
	if ci.Exponent == 0 {
		return fmt.Sprintf("%s%d %s", sign, u, m.Currency)
	}
	return fmt.Sprintf("%s%d.%0*d %s", sign, u/d, ci.Exponent, u%d, m.Currency)
}

// SumMoney sums all the monies in ms.
// All monies must be in the same currency c.
// If ms is empty, a valid zero money in currency c is returned.
func SumMoney(c Currency, ms []Money) (Money, error) {
	if !c.IsValid() {
		return Money{}, fmt.Errorf("%w: %v", ErrInvalidCurrency, c)
	}
	sum := big.NewInt(0)
	for _, m := range ms {
		if m.Currency != c {
			return Money{}, fmt.Errorf("%w: %q in a sum of %q", ErrCurrencyMismatch, m.Currency, c)
		}
		sum.Add(sum, big.NewInt(m.MinorUnits))
	}
	if !sum.IsInt64() {
		return Money{}, fmt.Errorf("%w: %v cannot be represented by int64", ErrOverflow, sum)
	}
	return Money{Currency: c, MinorUnits: sum.Int64()}, nil
}
