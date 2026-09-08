package saldo

import (
	"errors"
	"math"
	"testing"
)

func TestJournalEntryValidate(t *testing.T) {
	testCases := []struct {
		name    string
		fc      Currency
		je      JournalEntry
		wantErr error
	}{
		{
			name:    "no postings",
			fc:      USD,
			je:      JournalEntry{},
			wantErr: ErrNotEnoughPostings,
		},
		{
			name:    "empty postings",
			fc:      USD,
			je:      JournalEntry{Postings: []Posting{}},
			wantErr: ErrNotEnoughPostings,
		},
		{
			name: "one posting",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{0, USD}},
			}},
			wantErr: ErrNotEnoughPostings,
		},
		{
			name: "currency is checked before posting count",
			fc:   "XXX",
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{1, EUR}},
			}},
			wantErr: ErrInvalidCurrency,
		},
		{
			name: "two postings balance",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-500, USD}},
			}},
			wantErr: nil,
		},
		{
			name: "three postings balance",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-200, USD}},
				{FunctionalAmount: Money{-300, USD}},
			}},
			wantErr: nil,
		},
		{
			name: "all zero postings balance",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{0, USD}},
				{FunctionalAmount: Money{0, USD}},
			}},
			wantErr: nil,
		},
		{
			name: "currency without minor units",
			fc:   JPY,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{5, JPY}},
				{FunctionalAmount: Money{-5, JPY}},
			}},
			wantErr: nil,
		},
		{
			name: "fields other than Postings are ignored",
			fc:   USD,
			je: JournalEntry{
				JournalEntryID:            "je-1",
				LedgerID:                  "ledger-1",
				IdempotencyKey:            "key-1",
				SourceDocumentReferenceID: "invoice-1",
				Description:               "a sale",
				Reverses:                  "je-0",
				Postings: []Posting{
					{FunctionalAmount: Money{500, USD}},
					{FunctionalAmount: Money{-500, USD}},
				},
			},
			wantErr: nil,
		},
		{
			name: "both postings positive",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{300, USD}},
			}},
			wantErr: ErrPostingsSumNotZero,
		},
		{
			name: "both postings negative",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{-500, USD}},
				{FunctionalAmount: Money{-300, USD}},
			}},
			wantErr: ErrPostingsSumNotZero,
		},
		{
			name: "off by one minor unit",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-499, USD}},
			}},
			wantErr: ErrPostingsSumNotZero,
		},
		{
			name: "invalid functional currency",
			fc:   "XXX",
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-500, USD}},
			}},
			wantErr: ErrInvalidCurrency,
		},
		{
			// The zero value of Ledger.FunctionalCurrency.
			name: "empty functional currency",
			fc:   "",
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-500, USD}},
			}},
			wantErr: ErrInvalidCurrency,
		},
		{
			name: "posting in another currency",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-500, EUR}},
			}},
			wantErr: ErrCurrencyMismatch,
		},
		{
			name: "posting in an invalid currency",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-500, "XXX"}},
			}},
			wantErr: ErrCurrencyMismatch,
		},
		{
			name: "posting with no currency",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{500, USD}},
				{FunctionalAmount: Money{-500, ""}},
			}},
			wantErr: ErrCurrencyMismatch,
		},
		{
			name: "sum overflows int64",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{math.MaxInt64, USD}},
				{FunctionalAmount: Money{math.MaxInt64, USD}},
			}},
			wantErr: ErrOverflow,
		},
		{
			// A running int64 total would overflow on the second posting, but
			// the entry balances, so Validate must accept it.
			name: "partial sum overflows, entry still balances",
			fc:   USD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{math.MaxInt64, USD}},
				{FunctionalAmount: Money{math.MaxInt64, USD}},
				{FunctionalAmount: Money{math.MinInt64, USD}},
				{FunctionalAmount: Money{math.MinInt64, USD}},
				{FunctionalAmount: Money{2, USD}},
			}},
			wantErr: nil,
		},
		{
			// Validate balances FunctionalAmounts only. TransactionAmounts are
			// in whatever currency the event happened in and need not balance,
			// or even agree with each other.
			name: "transaction amounts are ignored",
			fc:   RSD,
			je: JournalEntry{Postings: []Posting{
				{FunctionalAmount: Money{11780, RSD}, TransactionAmount: Money{100, EUR}},
				{FunctionalAmount: Money{-11780, RSD}, TransactionAmount: Money{-7, JPY}},
			}},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.je.Validate(tc.fc); !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate(%q) with postings %v error = %v, want %v",
					tc.fc, tc.je.Postings, err, tc.wantErr)
			}
		})
	}
}
