package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type ListAccountsRequest struct{ LedgerID string }

type ListAccountsResponse struct{ Accounts []Account }

func (s *Service) ListAccounts(ctx context.Context, req ListAccountsRequest) (ListAccountsResponse, error) {
	// Validate input, create domain objects.
	ledgerID, err := domain.ParseLedgerID(req.LedgerID)
	if err != nil {
		return ListAccountsResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	// Get data from db.
	accounts, err := executeValue(ctx, s.uow, func(r Repositories) ([]domain.Account, error) {
		if _, err := r.Ledger.Get(ctx, ledgerID); err != nil {
			return nil, err
		}
		return r.Account.List(ctx, ledgerID)
	})

	// Map err.
	switch {
	case errors.Is(err, domain.ErrLedgerNotFound):
		return ListAccountsResponse{}, fmt.Errorf("%w: %w", ErrNotFound, err)
	case err != nil:
		return ListAccountsResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Map result.
	return ListAccountsResponse{accountsFromDomain(accounts)}, nil
}
