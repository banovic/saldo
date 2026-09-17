package domain

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

var (
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrInvalidMoney     = errors.New("invalid money")
	ErrOverflow         = errors.New("overflow")
)

// Money represents amount of money in some currency.
// Amount of money is expressed in minor units of that currency, because the focus of this package
// is to track movement of money (ie. not trading, crypto and other specific uses).
type Money struct {
	MinorUnits int64
	Currency   Currency
}

func (m Money) Validate() error {
	if err := m.Currency.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidMoney, err)
	}
	return nil
}

func (m Money) IsZeroAmount() bool {
	return m.MinorUnits == 0
}

func (m Money) Mul(k int64) (Money, error) {
	if err := m.Validate(); err != nil {
		return Money{}, err
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
