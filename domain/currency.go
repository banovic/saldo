package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
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

// Validates the currency. Currency is valid if has entry in currencyInfo map.
func (c Currency) Validate() error {
	if _, ok := c.Info(); !ok {
		return fmt.Errorf("%w: %q", ErrInvalidCurrency, c)
	}
	return nil
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
