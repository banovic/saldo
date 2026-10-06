package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidExchangeRateSource = errors.New("invalid exchange rate source")
)

type ExchangeRateSource string

const (
	NBS ExchangeRateSource = "nbs"
)

func (ers ExchangeRateSource) Validate() error {
	switch ers {
	case NBS:
		return nil
	}
	return fmt.Errorf("%w: %q", ErrInvalidExchangeRateSource, ers)
}
