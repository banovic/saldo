package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type CreateAccountRequest struct {
	LedgerID string
	Type     string
	Code     string
	Name     string
}

type CreateAccountResponse struct{ Account }

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (CreateAccountResponse, error) {
	// Create domain objects and validate them.
	ledgerID, err := domain.ParseLedgerID(req.LedgerID)
	if err != nil {
		return CreateAccountResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}
	account := domain.Account{
		AccountID: domain.AccountID{UUID: s.generateUUID()},
		LedgerID:  ledgerID,
		Type:      domain.AccountType(req.Type),
		Code:      domain.AccountCode(req.Code),
		Name:      domain.AccountName(req.Name),
		CreatedAt: s.utcTimestamp(),
	}
	if err := account.Validate(); err != nil {
		return CreateAccountResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	// Run db code inside unit of work:
	err = s.uow.Execute(ctx, func(r Repositories) error {
		// Make sure Ledger exists first.
		if _, err := r.Ledger.Get(ctx, ledgerID); err != nil {
			return err
		}
		// Insert only if Ledger is exists for given ledger id.
		return r.Account.Insert(ctx, account)
	})

	// Process errors.
	switch {
	case errors.Is(err, domain.ErrLedgerNotFound):
		return CreateAccountResponse{}, fmt.Errorf("%w: %w", ErrNotFound, err)
	case errors.Is(err, domain.ErrDuplicateAccountCodeInLedger):
		return CreateAccountResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	case err != nil:
		return CreateAccountResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Return domain objects mapped to DTO's.
	return CreateAccountResponse{accountFromDomain(account)}, nil
}
