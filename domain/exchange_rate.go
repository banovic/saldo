package domain

import (
	"errors"
	"fmt"
	"uuid"
)

var (
	ErrInvalidExchangeRate = errors.New("invalid exchange rate")
)

// ExchangeRateID is identifier for single ExchangeRate.
type ExchangeRateID struct{ uuid.UUID }

// IsZero checks if exchange rate id is zero.
func (erid ExchangeRateID) IsZero() bool {
	return erid.UUID == uuid.Nil()
}

// ExchangeRate is exchange rate between currency pair at certain point in time.
// Conversion must be deterministic - same inputs must produce same outputs always.
// ExchangeRate is append only, no deletions nor updates, once it is created.
// Exchange rate is saved as 2 integers representing: Num/Den.
// For example: exchange rate 1 EUR = 117.8 RSD is stored as Num: 1178, Den: 10
// Num/Den converts an amount in From Currency into an amount in To Currency.
// TODO!!! - this should be extended with exchange rate source (NBS, ECB ...), but
// later once there is need for it.
type ExchangeRate struct {
	ExchangeRateID ExchangeRateID
	From           Currency
	To             Currency
	Num            int64
	Den            int64
	On             Date
}

// Validate returns errors if exchange rate is not valid.
// ExchangeRate is valid if:
//   - exchange rate id is valid
//   - from currency is valid
//   - to currency is valid
//   - from and to currencies are different
//   - numerator must be greater than zero
//   - denominator must be greater than zero
//   - on must be valid date
func (er ExchangeRate) Validate() error {
	if er.ExchangeRateID.IsZero() {
		return fmt.Errorf("%w: exchange rate id is empty", ErrInvalidExchangeRate)
	}
	if err := er.From.Validate(); err != nil {
		return fmt.Errorf("%w: from: %w", ErrInvalidExchangeRate, err)
	}
	if err := er.To.Validate(); err != nil {
		return fmt.Errorf("%w: to: %w", ErrInvalidExchangeRate, err)
	}
	if er.From == er.To {
		return fmt.Errorf("%w: from and to must be different currencies: %q", ErrInvalidExchangeRate, er.From)
	}
	if er.Num <= 0 {
		return fmt.Errorf("%w: numerator must be greater than 0", ErrInvalidExchangeRate)
	}
	if er.Den <= 0 {
		return fmt.Errorf("%w: denominator must be greater than 0", ErrInvalidExchangeRate)
	}
	if err := er.On.Validate(); err != nil {
		return fmt.Errorf("%w: on: %w", ErrInvalidExchangeRate, err)
	}
	return nil
}
