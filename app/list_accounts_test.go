package app

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
)

func TestListAccounts(t *testing.T) {
	ledgerID := domain.LedgerID{UUID: uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")}
	cash := domain.Account{
		AccountID: domain.AccountID{UUID: uuid.MustParse("01992b6e-9a4b-7d2e-8c3f-5b6a7d8e9f0a")},
		LedgerID:  ledgerID,
		Type:      domain.Asset,
		Code:      "1",
		Name:      "Cash",
		CreatedAt: time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC),
	}
	sales := domain.Account{
		AccountID: domain.AccountID{UUID: uuid.MustParse("01992b6e-a55c-7e3f-9d4a-6c7b8e9f0a1b")},
		LedgerID:  ledgerID,
		Type:      domain.Revenue,
		Code:      "4",
		Name:      "Sales",
		CreatedAt: time.Date(2026, 9, 18, 11, 0, 0, 0, time.UTC),
	}
	cashDTO := Account{
		AccountID: cash.AccountID,
		LedgerID:  ledgerID,
		Type:      domain.Asset,
		Code:      "1",
		Name:      "Cash",
		CreatedAt: cash.CreatedAt,
	}
	salesDTO := Account{
		AccountID: sales.AccountID,
		LedgerID:  ledgerID,
		Type:      domain.Revenue,
		Code:      "4",
		Name:      "Sales",
		CreatedAt: sales.CreatedAt,
	}
	valid := ListAccountsRequest{LedgerID: ledgerID.String()}
	invalid := ListAccountsRequest{LedgerID: "not-a-uuid"}
	notFound := fmt.Errorf("%w: %v", domain.ErrLedgerNotFound, ledgerID)
	failure := errors.New("connection reset")

	testCases := []struct {
		name     string
		req      ListAccountsRequest
		ledgers  fakeLedgerRepository
		accounts fakeAccountRepository
		want     []Account
		wantErr  error
	}{
		{"keeps repository order", valid, fakeLedgerRepository{}, fakeAccountRepository{accounts: []domain.Account{sales, cash}}, []Account{salesDTO, cashDTO}, nil},
		{"no accounts", valid, fakeLedgerRepository{}, fakeAccountRepository{}, nil, nil},
		{"invalid ledger id", invalid, fakeLedgerRepository{}, fakeAccountRepository{}, nil, ErrInvalidInput},
		{"ledger not found", valid, fakeLedgerRepository{err: notFound}, fakeAccountRepository{}, nil, ErrNotFound},
		{"ledger not found keeps cause", valid, fakeLedgerRepository{err: notFound}, fakeAccountRepository{}, nil, domain.ErrLedgerNotFound},
		{"ledger repository failure", valid, fakeLedgerRepository{err: failure}, fakeAccountRepository{}, nil, ErrInternal},
		{"account repository failure", valid, fakeLedgerRepository{}, fakeAccountRepository{err: failure}, nil, ErrInternal},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(fakeUnitOfWork{ledgers: tc.ledgers, accounts: tc.accounts}, nil, nil)

			got, err := svc.ListAccounts(t.Context(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ListAccounts(%+v) error = %v, want %v", tc.req, err, tc.wantErr)
			}
			if !slices.Equal(got.Accounts, tc.want) {
				t.Errorf("ListAccounts(%+v).Accounts = %+v, want %+v", tc.req, got.Accounts, tc.want)
			}
		})
	}
}
