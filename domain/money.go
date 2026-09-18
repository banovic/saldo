package domain

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrOverflow         = errors.New("overflow")
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

// IsZeroAmount checks if money amount is zero.
func (m Money) IsZeroAmount() bool {
	return m.MinorUnits == 0
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
func (m Money) Mul(k int64) (Money, error) {
	if err := m.Validate(); err != nil {
		return Money{}, err
	}
	mul := new(big.Int).Mul(big.NewInt(k), big.NewInt(m.MinorUnits))
	if !mul.IsInt64() {
		return Money{}, fmt.Errorf("%w: %v cannot be represented by int64", ErrOverflow, mul)
	}
	return Money{MinorUnits: mul.Int64(), Currency: m.Currency}, nil
}

// SumMoney sums all the monies in ms.
// All monies must be in the same currency c.
// If ms is empty, a valid zero money in currency c is returned.
func SumMoney(c Currency, ms []Money) (Money, error) {
	if err := c.Validate(); err != nil {
		return Money{}, err
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
