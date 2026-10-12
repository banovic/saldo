package app

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
)

func TestCreateAccount(t *testing.T) {
	id := uuid.MustParse("01992b6e-9a4b-7d2e-8c3f-5b6a7d8e9f0a")
	ledgerID := domain.LedgerID{UUID: uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")}
	createdAt := time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC)
	now := func() time.Time { return createdAt }
	valid := CreateAccountRequest{LedgerID: ledgerID.String(), Type: "asset", Code: "1", Name: "Cash"}
	created := CreateAccountResponse{Account{
		AccountID: domain.AccountID{UUID: id},
		LedgerID:  ledgerID,
		Type:      domain.Asset,
		Code:      "1",
		Name:      "Cash",
		CreatedAt: createdAt,
	}}
	invalidID := CreateAccountRequest{LedgerID: "not-a-uuid", Type: "asset", Code: "1", Name: "Cash"}
	invalidType := CreateAccountRequest{LedgerID: ledgerID.String(), Type: "income", Code: "1", Name: "Cash"}
	notFound := fmt.Errorf("%w: %v", domain.ErrLedgerNotFound, ledgerID)
	duplicate := fmt.Errorf("%w: code: %q, ledger id: %s", domain.ErrDuplicateAccountCodeInLedger, valid.Code, ledgerID)

	testCases := []struct {
		name      string
		req       CreateAccountRequest
		getErr    error
		insertErr error
		want      CreateAccountResponse
		wantErr   error
	}{
		{"valid", valid, nil, nil, created, nil},
		{"invalid ledger id", invalidID, nil, nil, CreateAccountResponse{}, ErrInvalidInput},
		{"invalid ledger id keeps cause", invalidID, nil, nil, CreateAccountResponse{}, domain.ErrInvalidLedgerID},
		{"invalid account", invalidType, nil, nil, CreateAccountResponse{}, ErrInvalidInput},
		{"invalid account keeps cause", invalidType, nil, nil, CreateAccountResponse{}, domain.ErrInvalidAccount},
		{"ledger not found", valid, notFound, nil, CreateAccountResponse{}, ErrNotFound},
		{"ledger not found keeps cause", valid, notFound, nil, CreateAccountResponse{}, domain.ErrLedgerNotFound},
		{"duplicate code", valid, nil, duplicate, CreateAccountResponse{}, ErrInvalidInput},
		{"duplicate code keeps cause", valid, nil, duplicate, CreateAccountResponse{}, domain.ErrDuplicateAccountCodeInLedger},
		{"ledger repository failure", valid, errors.New("connection reset"), nil, CreateAccountResponse{}, ErrInternal},
		{"account repository failure", valid, nil, errors.New("connection reset"), CreateAccountResponse{}, ErrInternal},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uow := fakeUnitOfWork{
				ledgers:  fakeLedgerRepository{err: tc.getErr},
				accounts: fakeAccountRepository{err: tc.insertErr},
			}
			svc := NewService(uow, now, func() uuid.UUID { return id })

			got, err := svc.CreateAccount(t.Context(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("CreateAccount(%+v) error = %v, want %v", tc.req, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("CreateAccount(%+v) = %+v, want %+v", tc.req, got, tc.want)
			}
		})
	}
}
