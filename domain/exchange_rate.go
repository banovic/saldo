package domain

import "time"

// ExchangeRateID is ID of certain ExchangeRate, it can be empty ("" is valid).
type ExchangeRateID string

// ExchangeRateKind says what a rate is for. Different accounting
// treatments require different rates for the same pair and date.
type ExchangeRateKind string

const (
	Spot    ExchangeRateKind = "spot"    // transaction-date rate, for initial recognition
	Closing ExchangeRateKind = "closing" // period-end rate, for retranslating monetary items
	Average ExchangeRateKind = "average" // period average, for translating income and expenses
)

// ExchangeRateSource is the source of given ExchangeRate
type ExchangeRateSource string

const (
	NBS       ExchangeRateSource = "nbs"       // Narodna Banka Srbije.
	ECB       ExchangeRateSource = "ecb"       // European Central Bank.
	Processor ExchangeRateSource = "processor" // The rate the PSP actually applied
	Contract  ExchangeRateSource = "contract"  // Fixed by agreement.
	Manual    ExchangeRateSource = "manual"    // Entered by hand.
)

// ExchangeRate is exchange rate between currency pair at certain point in time
// and from certain source.
// ExchangeRate is append only, no deletions nor updates, once it is created.
// Exchange rate is saved as 2 integers representing: Num/Den.
// For example: exchange rate 1 EUR = 117.8 RSD is stored as Num: 1178, Den: 10
type ExchangeRate struct {
	ExchangeRateID ExchangeRateID
	From           Currency
	To             Currency

	// Numerator
	Num uint64
	// Denominator
	Den uint64

	Source      ExchangeRateSource
	Kind        ExchangeRateKind
	On          time.Time
	Description string
}
