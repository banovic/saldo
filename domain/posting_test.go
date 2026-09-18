package domain

import (
	"errors"
	"math"
	"testing"
)

func TestPostingValidate(t *testing.T) {
	eurRsd := ExchangeRate{ExchangeRateID: "r1", From: EUR, To: RSD, Num: 1178, Den: 10, Source: NBS, Kind: Spot}

	testCases := []struct {
		name    string
		p       Posting
		er      ExchangeRate
		wantErr error
	}{
		{
			name: "same currency, no rate",
			p:    Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{500, USD}},
		},
		{
			name: "same currency, no rate, negative",
			p:    Posting{TransactionAmount: Money{-500, USD}, FunctionalAmount: Money{-500, USD}},
		},
		{
			name: "different currencies, rate set",
			p:    Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{11780, RSD}},
			er:   eurRsd,
		},
		{
			name: "different currencies, rate set, negative",
			p:    Posting{TransactionAmount: Money{-100, EUR}, FunctionalAmount: Money{-11780, RSD}},
			er:   eurRsd,
		},
		{
			name:    "invalid transaction currency",
			p:       Posting{TransactionAmount: Money{500, "XXX"}, FunctionalAmount: Money{500, USD}},
			wantErr: ErrInvalidCurrency,
		},
		{
			name:    "zero value transaction amount",
			p:       Posting{FunctionalAmount: Money{500, USD}},
			wantErr: ErrInvalidCurrency,
		},
		{
			name:    "invalid functional currency",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{500, "XXX"}},
			wantErr: ErrInvalidCurrency,
		},
		{
			name:    "sign mismatch, positive transaction",
			p:       Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{-11780, RSD}},
			er:      eurRsd,
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "sign mismatch, negative transaction",
			p:       Posting{TransactionAmount: Money{-100, EUR}, FunctionalAmount: Money{11780, RSD}},
			er:      eurRsd,
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "sign mismatch, zero functional",
			p:       Posting{TransactionAmount: Money{1, EUR}, FunctionalAmount: Money{0, RSD}},
			er:      eurRsd,
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "same currency, no rate, amounts differ",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{501, USD}},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "different currencies, no rate",
			p:       Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{11780, RSD}},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "different currencies, same minor units, no rate",
			p:       Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{100, USD}},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "same currency, rate set",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{500, USD}},
			er:      eurRsd,
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "same currency, rate set, amounts differ",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{501, USD}},
			er:      eurRsd,
			wantErr: ErrInvalidPosting,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.p.Validate(tc.er)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%+v.Validate(%+v) = %v, want %v", tc.p, tc.er, err, tc.wantErr)
			}
			if tc.wantErr != nil && !errors.Is(err, ErrInvalidPosting) {
				t.Errorf("%+v.Validate(%+v) = %v, want %v", tc.p, tc.er, err, ErrInvalidPosting)
			}
		})
	}
}

func TestPostingIsDebitIsCredit(t *testing.T) {
	testCases := []struct {
		name       string
		fa         Money
		wantDebit  bool
		wantCredit bool
	}{
		{"positive is debit", Money{1, USD}, true, false},
		{"negative is credit", Money{-1, USD}, false, true},
		{"zero is neither", Money{0, USD}, false, false},
		{"MaxInt64 is debit", Money{math.MaxInt64, USD}, true, false},
		{"MinInt64 is credit", Money{math.MinInt64, USD}, false, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Posting{FunctionalAmount: tc.fa}
			if got := p.IsDebit(); got != tc.wantDebit {
				t.Errorf("Posting{FunctionalAmount: %v}.IsDebit() = %t, want %t", tc.fa, got, tc.wantDebit)
			}
			if got := p.IsCredit(); got != tc.wantCredit {
				t.Errorf("Posting{FunctionalAmount: %v}.IsCredit() = %t, want %t", tc.fa, got, tc.wantCredit)
			}
		})
	}
}
