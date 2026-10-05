package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
)

type fakeUnitOfWork struct {
	ledgers domain.LedgerRepository
}

func (f fakeUnitOfWork) Execute(ctx context.Context, work func(Repositories) error) error {
	return work(Repositories{Ledger: f.ledgers})
}

// fakeLedgerRepository embeds the interface so only methods under test need implementing.
type fakeLedgerRepository struct {
	domain.LedgerRepository
	insertErr error
}

func (f fakeLedgerRepository) Insert(ctx context.Context, l domain.Ledger) error {
	return f.insertErr
}

func TestCreateLedger(t *testing.T) {
	id := uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")
	now := func() time.Time { return time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC) }
	valid := CreateLedgerRequest{Name: "Acme", FunctionalCurrency: "USD", ReportingTimeZone: "Europe/Belgrade"}
	duplicate := fmt.Errorf("%w: %q", domain.ErrDuplicateLedgerName, valid.Name)

	testCases := []struct {
		name      string
		req       CreateLedgerRequest
		insertErr error
		wantErr   error
	}{
		{"valid", valid, nil, nil},
		{"invalid request", CreateLedgerRequest{Name: "ab", FunctionalCurrency: "USD", ReportingTimeZone: "Europe/Belgrade"}, nil, ErrInvalidInput},
		{"duplicate name", valid, duplicate, ErrInvalidInput},
		{"duplicate name keeps cause", valid, duplicate, domain.ErrDuplicateLedgerName},
		{"repository failure", valid, errors.New("connection reset"), ErrInternal},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uow := fakeUnitOfWork{ledgers: fakeLedgerRepository{insertErr: tc.insertErr}}
			svc := NewService(uow, now, func() uuid.UUID { return id })

			resp, err := svc.CreateLedger(t.Context(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("CreateLedger(%+v) error = %v, want %v", tc.req, err, tc.wantErr)
			}
			if want := (domain.LedgerID{UUID: id}); tc.wantErr == nil && resp.LedgerID != want {
				t.Errorf("CreateLedger(%+v).LedgerID = %v, want %v", tc.req, resp.LedgerID, want)
			}
		})
	}
}
