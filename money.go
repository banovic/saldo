package saldo

import (
	"errors"
	"fmt"
	"math"
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

// Returns (CurrencyInfo, true) for called Currency, if Currency is valid and defined.
// Returns (zero CurrencyInfo, false) if called Currency is invalid.
// Info() returns the CurrencyInfo for c and reports whether c is a defined currency.
func (c Currency) Info() (CurrencyInfo, bool) {
	ci, ok := currencyInfo[c]
	return ci, ok
}

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

func NewMoney(units int64, currency Currency) (Money, error) {
	if !currency.IsValid() {
		return Money{}, fmt.Errorf("%w: %q", ErrInvalidCurrency, currency)
	}
	return Money{MinorUnits: units, Currency: currency}, nil
}

func ZeroMoney(currency Currency) (Money, error) {
	return NewMoney(0, currency)
}

func (a Money) IsValid() bool {
	return a.Currency.IsValid()
}

func (a Money) Mul(k int64) (Money, error) {
	if !a.IsValid() {
		return Money{}, fmt.Errorf("%w: %v", ErrInvalidMoney, a)
	}
	if a.MinorUnits == 0 || k == 0 {
		return Money{MinorUnits: 0, Currency: a.Currency}, nil
	}
	// p/k test below cannot detect these values; see 'Go defines MinInt64 / -1 as MinInt64'
	// math.MinInt64 * -1 = math.MinInt64, but this value has overflown already
	// math.MinInt64 / -1 = math.MinInt64, but this value has overflown already
	if a.MinorUnits == math.MinInt64 && k == -1 {
		return Money{}, fmt.Errorf("%w: %v times %d", ErrOverflow, a, k)
	}
	p := a.MinorUnits * k
	if p/k != a.MinorUnits {
		return Money{}, fmt.Errorf("%w: %v times %d", ErrOverflow, a, k)
	}
	return Money{MinorUnits: p, Currency: a.Currency}, nil
}

// Displaying money amounts to end user usually takes into account user's location.
// This method does not do that - it's just for simple display of Money struct.
func (a Money) String() string {
	ci, ok1 := a.Currency.Info()
	d, ok2 := ci.MinorUnitsPerUnit()
	if !ok1 || !ok2 {
		return fmt.Sprintf("%d minor units of %q", a.MinorUnits, a.Currency)
	}
	sign, u := "", uint64(a.MinorUnits)
	if a.MinorUnits < 0 {
		sign, u = "-", -u
	}
	if ci.Exponent == 0 {
		return fmt.Sprintf("%s%d %s", sign, u, a.Currency)
	}
	return fmt.Sprintf("%s%d.%0*d %s", sign, u/d, ci.Exponent, u%d, a.Currency)
}
