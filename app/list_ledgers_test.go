package app

import (
	"errors"
	"slices"
	"testing"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
)

func TestListLedgers(t *testing.T) {
	acme := domain.Ledger{
		LedgerID:           domain.LedgerID{UUID: uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")},
		Name:               "Acme",
		FunctionalCurrency: domain.USD,
		ReportingTimeZone:  "Europe/Belgrade",
		CreatedAt:          time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC),
	}
	globex := domain.Ledger{
		LedgerID:           domain.LedgerID{UUID: uuid.MustParse("01992b6e-9a4b-7d2e-8c3f-5b6a7d8e9f0a")},
		Name:               "Globex",
		FunctionalCurrency: domain.RSD,
		ReportingTimeZone:  "UTC",
		CreatedAt:          time.Date(2026, 9, 18, 11, 0, 0, 0, time.UTC),
	}
	acmeDTO := Ledger{
		LedgerID:           acme.LedgerID,
		Name:               "Acme",
		FunctionalCurrency: domain.USD,
		ReportingTimeZone:  "Europe/Belgrade",
		CreatedAt:          acme.CreatedAt,
	}
	globexDTO := Ledger{
		LedgerID:           globex.LedgerID,
		Name:               "Globex",
		FunctionalCurrency: domain.RSD,
		ReportingTimeZone:  "UTC",
		CreatedAt:          globex.CreatedAt,
	}

	testCases := []struct {
		name    string
		repo    fakeLedgerRepository
		want    []Ledger
		wantErr error
	}{
		{"keeps repository order", fakeLedgerRepository{ledgers: []domain.Ledger{globex, acme}}, []Ledger{globexDTO, acmeDTO}, nil},
		{"no ledgers", fakeLedgerRepository{}, nil, nil},
		{"repository failure", fakeLedgerRepository{err: errors.New("connection reset")}, nil, ErrInternal},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(fakeUnitOfWork{ledgers: tc.repo}, nil, nil)

			got, err := svc.ListLedgers(t.Context(), ListLedgersRequest{})
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ListLedgers() error = %v, want %v", err, tc.wantErr)
			}
			if !slices.Equal(got.Ledgers, tc.want) {
				t.Errorf("ListLedgers().Ledgers = %+v, want %+v", got.Ledgers, tc.want)
			}
		})
	}
}
