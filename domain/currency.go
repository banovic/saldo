package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
)

// Currency is ISO 4217 3 letter code.
type Currency string

// CurrencyInfo holds other ISO 4217 information about currency.
type CurrencyInfo struct {
	Name             string
	Num              string
	Exponent         int
	MinorUnitsInUnit int64
}

const (
	EUR Currency = "EUR"
	RSD Currency = "RSD"
	USD Currency = "USD"
	JPY Currency = "JPY"
	TND Currency = "TND"
)

var currencyInfo = map[Currency]CurrencyInfo{
	EUR: {Name: "Euro", Num: "978", Exponent: 2, MinorUnitsInUnit: 100},
	RSD: {Name: "Serbian dinar", Num: "941", Exponent: 2, MinorUnitsInUnit: 100},
	USD: {Name: "US Dollar", Num: "840", Exponent: 2, MinorUnitsInUnit: 100},
	JPY: {Name: "Japanese yen", Num: "392", Exponent: 0, MinorUnitsInUnit: 1},
	TND: {Name: "Tunisian dinar", Num: "788", Exponent: 3, MinorUnitsInUnit: 1000},
}

// Info returns the CurrencyInfo for c and reports whether c is a defined currency.
func (c Currency) Info() (CurrencyInfo, bool) {
	ci, ok := currencyInfo[c]
	return ci, ok
}

// Validate returns error if c is not defined currency.
func (c Currency) Validate() error {
	if _, ok := c.Info(); !ok {
		return fmt.Errorf("%w: %q", ErrInvalidCurrency, c)
	}
	return nil
}
