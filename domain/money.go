package domain

import (
	"errors"
	"fmt"
)

var (
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrInvalidMoney     = errors.New("invalid money")
)

// Money represents amount of money in some currency.
// Amount of money is expressed in minor units of that currency, because the focus of this package
// is to track movement of money (ie. not trading, crypto and other specific uses).
type Money struct {
	MinorUnits int64
	Currency   Currency
}

// Validate returns errors if money currency is not valid.
func (m Money) Validate() error {
	if err := m.Currency.Validate(); err != nil {
		return fmt.Errorf("money: %w", err)
	}
	return nil
}

// IsZero checks if money amount is zero.
func (m Money) IsZero() bool {
	return m.MinorUnits == 0
}

// Sign returns 1, 0, -1 if minor units is greater, equal or less than zero.
func (m Money) Sign() int {
	switch {
	case m.MinorUnits < 0:
		return -1
	case m.MinorUnits > 0:
		return 1
	default:
		return 0
	}
}

// String returns a simple representation of money m.
// Displaying money amounts to end user usually takes into account user's location.
// This method does not do that - it's just for simple display of Money struct.
func (m Money) String() string {
	ci, ok := m.Currency.Info()
	if !ok {
		return fmt.Sprintf("%d minor units of %q", m.MinorUnits, m.Currency)
	}
	// Negate as uint64; -MinInt64 does not fit in int64
	sign, u := "", uint64(m.MinorUnits)
	if m.MinorUnits < 0 {
		sign, u = "-", -u
	}
	d := uint64(ci.MinorUnitsInUnit)
	if d <= 1 {
		return fmt.Sprintf("%s%d %s", sign, u, m.Currency)
	}
	return fmt.Sprintf("%s%d.%0*d %s", sign, u/d, ci.Exponent, u%d, m.Currency)
}

// Mul multiplies money by given factor k.
// Money instance is assumed to be valid when calling this method.
func (m Money) Mul(k int64) (Money, error) {
	newUnits, err := mul(m.MinorUnits, k)
	if err != nil {
		return Money{}, fmt.Errorf("%w: %w", ErrInvalidMoney, err)
	}
	return Money{MinorUnits: newUnits, Currency: m.Currency}, nil
}

// Sum sums all the monies in ms.
// All monies must be in the same currency c.
// If ms is empty, a valid zero money in currency c is returned.
// Money instances are assumed to be valid when calling this method.
// Currency c is assumed to be valid when calling this method.
func Sum(c Currency, ms []Money) (Money, error) {
	monies := make([]int64, 0, len(ms))
	for _, m := range ms {
		if m.Currency != c {
			return Money{}, fmt.Errorf("%w: %q in a sum of %q", ErrCurrencyMismatch, m.Currency, c)
		}
		monies = append(monies, m.MinorUnits)
	}
	moniesSum, err := add(monies...)
	if err != nil {
		return Money{}, fmt.Errorf("%w: %w", ErrInvalidMoney, err)
	}
	return Money{MinorUnits: moniesSum, Currency: c}, nil
}
