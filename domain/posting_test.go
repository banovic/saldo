package domain

import (
	"errors"
	"math"
	"testing"
	"uuid"
)

func TestPostingValidate(t *testing.T) {
	rateID := ExchangeRateID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8eb0")}

	testCases := []struct {
		name    string
		p       Posting
		wantErr error
	}{
		{
			name: "same currency, no rate id",
			p:    Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{500, USD}},
		},
		{
			name: "same currency, no rate id, negative",
			p:    Posting{TransactionAmount: Money{-500, USD}, FunctionalAmount: Money{-500, USD}},
		},
		{
			name: "different currencies, rate id set",
			p:    Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{11780, RSD}, ExchangeRateID: rateID},
		},
		{
			name: "different currencies, rate id set, negative",
			p:    Posting{TransactionAmount: Money{-100, EUR}, FunctionalAmount: Money{-11780, RSD}, ExchangeRateID: rateID},
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
			p:       Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{-11780, RSD}, ExchangeRateID: rateID},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "sign mismatch, negative transaction",
			p:       Posting{TransactionAmount: Money{-100, EUR}, FunctionalAmount: Money{11780, RSD}, ExchangeRateID: rateID},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "sign mismatch, zero functional",
			p:       Posting{TransactionAmount: Money{1, EUR}, FunctionalAmount: Money{0, RSD}, ExchangeRateID: rateID},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "same currency, no rate, amounts differ",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{501, USD}},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "different currencies, no rate id",
			p:       Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{11780, RSD}},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "different currencies, same minor units, no rate id",
			p:       Posting{TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{100, USD}},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "same currency, rate id set",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{500, USD}, ExchangeRateID: rateID},
			wantErr: ErrInvalidPosting,
		},
		{
			name:    "same currency, rate id set, amounts differ",
			p:       Posting{TransactionAmount: Money{500, USD}, FunctionalAmount: Money{501, USD}, ExchangeRateID: rateID},
			wantErr: ErrInvalidPosting,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.p.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%+v.Validate() = %v, want %v", tc.p, err, tc.wantErr)
			}
			if tc.wantErr != nil && !errors.Is(err, ErrInvalidPosting) {
				t.Errorf("%+v.Validate() = %v, want %v", tc.p, err, ErrInvalidPosting)
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
