package domain

import (
	"errors"
	"testing"
	"uuid"
)

func TestAccountIDIsZero(t *testing.T) {
	testCases := []struct {
		name string
		id   AccountID
		want bool
	}{
		{"zero value", AccountID{}, true},
		{"nil uuid", AccountID{uuid.Nil()}, true},
		{"non-nil uuid", AccountID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ed0")}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.IsZero(); got != tc.want {
				t.Errorf("%v.IsZero() = %t, want %t", tc.id, got, tc.want)
			}
		})
	}
}

func TestAccountValidate(t *testing.T) {
	valid := Account{
		AccountID: AccountID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8ed0")},
		LedgerID:  LedgerID{uuid.MustParse("01992b6e-8f3a-7c1d-9b2e-4a5f6c7d8e9f")},
		Type:      Asset,
		Code:      "241",
		Name:      "Tekući račun",
	}
	with := func(f func(*Account)) Account {
		a := valid
		f(&a)
		return a
	}

	testCases := []struct {
		name    string
		account Account
		wantErr error
	}{
		{"valid", valid, nil},
		{"zero account id", with(func(a *Account) { a.AccountID = AccountID{} }), ErrInvalidAccount},
		{"zero ledger id", with(func(a *Account) { a.LedgerID = LedgerID{} }), ErrInvalidAccount},
		{"invalid type", with(func(a *Account) { a.Type = "income" }), ErrInvalidAccountType},
		{"empty type", with(func(a *Account) { a.Type = "" }), ErrInvalidAccountType},
		{"empty code", with(func(a *Account) { a.Code = "" }), ErrInvalidAccountCode},
		{"invalid name", with(func(a *Account) { a.Name = "ab" }), ErrInvalidAccountName},
		{"empty name", with(func(a *Account) { a.Name = "" }), ErrInvalidAccountName},
		{"type is checked before code", with(func(a *Account) {
			a.Type = ""
			a.Code = ""
		}), ErrInvalidAccountType},
		{"code is checked before name", with(func(a *Account) {
			a.Code = ""
			a.Name = ""
		}), ErrInvalidAccountCode},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.account.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%+v.Validate() = %v, want %v", tc.account, err, tc.wantErr)
			}
			if tc.wantErr != nil && !errors.Is(err, ErrInvalidAccount) {
				t.Errorf("%+v.Validate() = %v, want %v", tc.account, err, ErrInvalidAccount)
			}
		})
	}
}
