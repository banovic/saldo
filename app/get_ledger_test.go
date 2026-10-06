package app

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
)

func TestGetLedger(t *testing.T) {
	ledger := domain.Ledger{
		LedgerID:           domain.LedgerID{UUID: uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")},
		Name:               "Acme",
		FunctionalCurrency: domain.USD,
		ReportingTimeZone:  "Europe/Belgrade",
		CreatedAt:          time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC),
	}
	found := GetLedgerResponse{Ledger{
		LedgerID:           ledger.LedgerID,
		Name:               "Acme",
		FunctionalCurrency: domain.USD,
		ReportingTimeZone:  "Europe/Belgrade",
		CreatedAt:          ledger.CreatedAt,
	}}
	valid := GetLedgerRequest{LedgerID: ledger.LedgerID.String()}
	invalid := GetLedgerRequest{LedgerID: "not-a-uuid"}
	notFound := fmt.Errorf("%w: %v", domain.ErrLedgerNotFound, ledger.LedgerID)

	testCases := []struct {
		name    string
		req     GetLedgerRequest
		repo    fakeLedgerRepository
		want    GetLedgerResponse
		wantErr error
	}{
		{"found", valid, fakeLedgerRepository{ledger: ledger}, found, nil},
		{"invalid id", invalid, fakeLedgerRepository{}, GetLedgerResponse{}, ErrInvalidInput},
		{"invalid id keeps cause", invalid, fakeLedgerRepository{}, GetLedgerResponse{}, domain.ErrInvalidLedgerID},
		{"not found", valid, fakeLedgerRepository{err: notFound}, GetLedgerResponse{}, ErrNotFound},
		{"not found keeps cause", valid, fakeLedgerRepository{err: notFound}, GetLedgerResponse{}, domain.ErrLedgerNotFound},
		{"repository failure", valid, fakeLedgerRepository{err: errors.New("connection reset")}, GetLedgerResponse{}, ErrInternal},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(fakeUnitOfWork{ledgers: tc.repo}, nil, nil)

			got, err := svc.GetLedger(t.Context(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("GetLedger(%+v) error = %v, want %v", tc.req, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("GetLedger(%+v) = %+v, want %+v", tc.req, got, tc.want)
			}
		})
	}
}
