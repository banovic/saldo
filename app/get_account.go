package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type GetAccountRequest struct {
	AccountID string
}

type GetAccountResponse struct{ Account }

func (s *Service) GetAccount(ctx context.Context, req GetAccountRequest) (GetAccountResponse, error) {
	// Create domain objects and validate.
	accountID, err := domain.ParseAccountID(req.AccountID)
	if err != nil {
		return GetAccountResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	// Get values from db.
	account, err := executeValue(ctx, s.uow, func(r Repositories) (domain.Account, error) {
		return r.Account.Get(ctx, accountID)
	})

	// Translate errors.
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		return GetAccountResponse{}, fmt.Errorf("%w: %w", ErrNotFound, err)
	case err != nil:
		return GetAccountResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Map domain object to DTO and return result.
	return GetAccountResponse{accountFromDomain(account)}, nil
}
