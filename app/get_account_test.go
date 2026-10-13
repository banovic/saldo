package app

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
)

func TestGetAccount(t *testing.T) {
	account := domain.Account{
		AccountID: domain.AccountID{UUID: uuid.MustParse("01992b6e-9a4b-7d2e-8c3f-5b6a7d8e9f0a")},
		LedgerID:  domain.LedgerID{UUID: uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")},
		Type:      domain.Asset,
		Code:      "1",
		Name:      "Cash",
		CreatedAt: time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC),
	}
	found := GetAccountResponse{Account{
		AccountID: account.AccountID,
		LedgerID:  account.LedgerID,
		Type:      domain.Asset,
		Code:      "1",
		Name:      "Cash",
		CreatedAt: account.CreatedAt,
	}}
	valid := GetAccountRequest{AccountID: account.AccountID.String()}
	invalid := GetAccountRequest{AccountID: "not-a-uuid"}
	notFound := fmt.Errorf("%w: %v", domain.ErrAccountNotFound, account.AccountID)

	testCases := []struct {
		name    string
		req     GetAccountRequest
		repo    fakeAccountRepository
		want    GetAccountResponse
		wantErr error
	}{
		{"found", valid, fakeAccountRepository{account: account}, found, nil},
		{"invalid id", invalid, fakeAccountRepository{}, GetAccountResponse{}, ErrInvalidInput},
		{"not found", valid, fakeAccountRepository{err: notFound}, GetAccountResponse{}, ErrNotFound},
		{"not found keeps cause", valid, fakeAccountRepository{err: notFound}, GetAccountResponse{}, domain.ErrAccountNotFound},
		{"repository failure", valid, fakeAccountRepository{err: errors.New("connection reset")}, GetAccountResponse{}, ErrInternal},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(fakeUnitOfWork{accounts: tc.repo}, nil, nil)

			got, err := svc.GetAccount(t.Context(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("GetAccount(%+v) error = %v, want %v", tc.req, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("GetAccount(%+v) = %+v, want %+v", tc.req, got, tc.want)
			}
		})
	}
}
