package domain

import (
	"errors"
	"math"
	"testing"
	"time"
	"uuid"
)

func TestJournalEntryIDIsZero(t *testing.T) {
	testCases := []struct {
		name string
		id   JournalEntryID
		want bool
	}{
		{"zero value", JournalEntryID{}, true},
		{"nil uuid", JournalEntryID{uuid.Nil()}, true},
		{"non-nil uuid", JournalEntryID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ea0")}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.IsZero(); got != tc.want {
				t.Errorf("%v.IsZero() = %t, want %t", tc.id, got, tc.want)
			}
		})
	}
}

func TestJournalEntryValidate(t *testing.T) {
	var (
		ledgerID      = LedgerID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")}
		otherLedgerID = LedgerID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9e")}
		entryID       = JournalEntryID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ea0")}
		otherEntryID  = JournalEntryID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ea1")}
		rateID        = ExchangeRateID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8eb0")}
		postingID     = PostingID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ec0")}
	)

	usd := func(n int64) Posting {
		return Posting{PostingID: postingID, TransactionAmount: Money{n, USD}, FunctionalAmount: Money{n, USD}}
	}

	// 2026-09-18 23:30 UTC is 2026-09-19 01:30 in Europe/Belgrade (UTC+2).
	occurredAt := time.Date(2026, time.September, 18, 23, 30, 0, 0, time.UTC)

	// validLedger and validEntry return fresh values; each test case edits its own copy.
	validLedger := func() Ledger {
		return Ledger{
			LedgerID:           ledgerID,
			Name:               "Books",
			FunctionalCurrency: USD,
			ReportingTimeZone:  "Europe/Belgrade",
		}
	}
	validEntry := func() JournalEntry {
		return JournalEntry{
			JournalEntryID: entryID,
			LedgerID:       ledgerID,
			IdempotencyKey: "key-1",
			OccurredAt:     occurredAt,
			PostedOn:       Date{2026, time.September, 19},
			Postings:       []Posting{usd(500), usd(-500)},
		}
	}

	testCases := []struct {
		name    string
		edit    func(je *JournalEntry, l *Ledger)
		wantErr error
	}{
		{
			name: "valid",
			edit: func(je *JournalEntry, l *Ledger) {},
		},
		{
			name: "three postings balance",
			edit: func(je *JournalEntry, l *Ledger) {
				je.Postings = []Posting{usd(500), usd(-200), usd(-300)}
			},
		},
		{
			name: "currency without minor units",
			edit: func(je *JournalEntry, l *Ledger) {
				l.FunctionalCurrency = JPY
				je.Postings = []Posting{
					{PostingID: postingID, TransactionAmount: Money{5, JPY}, FunctionalAmount: Money{5, JPY}},
					{PostingID: postingID, TransactionAmount: Money{-5, JPY}, FunctionalAmount: Money{-5, JPY}},
				}
			},
		},
		{
			// Only FunctionalAmounts must balance. TransactionAmounts are
			// in whatever currency the event happened in.
			name: "transaction amounts in other currencies are not balanced",
			edit: func(je *JournalEntry, l *Ledger) {
				l.FunctionalCurrency = RSD
				je.Postings = []Posting{
					{PostingID: postingID, TransactionAmount: Money{100, EUR}, FunctionalAmount: Money{11780, RSD}, ExchangeRateID: rateID},
					{PostingID: postingID, TransactionAmount: Money{-7, JPY}, FunctionalAmount: Money{-11780, RSD}, ExchangeRateID: rateID},
				}
			},
		},
		{
			// A running int64 total would overflow on the second posting, but
			// the entry balances, so Validate must accept it.
			name: "partial sum overflows, entry still balances",
			edit: func(je *JournalEntry, l *Ledger) {
				je.Postings = []Posting{usd(math.MaxInt64), usd(math.MaxInt64), usd(math.MinInt64), usd(math.MinInt64), usd(2)}
			},
		},
		{
			name: "reverses another entry",
			edit: func(je *JournalEntry, l *Ledger) { je.Reverses = otherEntryID },
		},
		{
			name: "posted after occurred",
			edit: func(je *JournalEntry, l *Ledger) { je.PostedOn = Date{2026, time.October, 1} },
		},
		{
			name: "posted on occurred date in UTC ledger",
			edit: func(je *JournalEntry, l *Ledger) {
				l.ReportingTimeZone = "UTC"
				je.PostedOn = Date{2026, time.September, 18}
			},
		},

		{
			name:    "zero journal entry id",
			edit:    func(je *JournalEntry, l *Ledger) { je.JournalEntryID = JournalEntryID{} },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "ledger id mismatch",
			edit:    func(je *JournalEntry, l *Ledger) { je.LedgerID = otherLedgerID },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "no postings",
			edit:    func(je *JournalEntry, l *Ledger) { je.Postings = nil },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "one posting",
			edit:    func(je *JournalEntry, l *Ledger) { je.Postings = []Posting{usd(0)} },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name: "invalid posting",
			edit: func(je *JournalEntry, l *Ledger) {
				je.Postings[0].TransactionAmount = Money{501, USD}
			},
			wantErr: ErrInvalidPosting,
		},
		{
			name: "posting not in functional currency",
			edit: func(je *JournalEntry, l *Ledger) {
				je.Postings[1] = Posting{PostingID: postingID, TransactionAmount: Money{-500, EUR}, FunctionalAmount: Money{-500, EUR}}
			},
			wantErr: ErrCurrencyMismatch,
		},
		{
			name:    "both postings positive",
			edit:    func(je *JournalEntry, l *Ledger) { je.Postings = []Posting{usd(500), usd(300)} },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "off by one minor unit",
			edit:    func(je *JournalEntry, l *Ledger) { je.Postings = []Posting{usd(500), usd(-499)} },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name: "sum overflows int64",
			edit: func(je *JournalEntry, l *Ledger) {
				je.Postings = []Posting{usd(math.MaxInt64), usd(math.MaxInt64)}
			},
			wantErr: ErrOverflow,
		},
		{
			name:    "reverses itself",
			edit:    func(je *JournalEntry, l *Ledger) { je.Reverses = entryID },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "empty idempotency key",
			edit:    func(je *JournalEntry, l *Ledger) { je.IdempotencyKey = "" },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "invalid ledger time zone",
			edit:    func(je *JournalEntry, l *Ledger) { l.ReportingTimeZone = "Europe/Nowhere" },
			wantErr: ErrInvalidTimeZone,
		},
		{
			// The UTC date of OccurredAt is 2026-09-18, but in the ledger's
			// time zone it is already 2026-09-19.
			name:    "posted before occurred in ledger time zone",
			edit:    func(je *JournalEntry, l *Ledger) { je.PostedOn = Date{2026, time.September, 18} },
			wantErr: ErrInvalidJournalEntry,
		},
		{
			name:    "zero posted on",
			edit:    func(je *JournalEntry, l *Ledger) { je.PostedOn = Date{} },
			wantErr: ErrInvalidJournalEntry,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			je, l := validEntry(), validLedger()
			tc.edit(&je, &l)
			err := je.Validate(l)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil && !errors.Is(err, ErrInvalidJournalEntry) {
				t.Errorf("Validate() error = %v, want %v", err, ErrInvalidJournalEntry)
			}
		})
	}
}
