package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type GetLedgerRequest struct {
	LedgerID string
}

type GetLedgerResponse struct{ Ledger }

func (s *Service) GetLedger(ctx context.Context, req GetLedgerRequest) (GetLedgerResponse, error) {
	ledgerID, err := domain.ParseLedgerID(req.LedgerID)
	if err != nil {
		return GetLedgerResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	// Get values from db.
	ledger, err := executeValue(ctx, s.uow, func(rs Repositories) (domain.Ledger, error) {
		return rs.Ledger.Get(ctx, ledgerID)
	})

	switch {
	case errors.Is(err, domain.ErrLedgerNotFound):
		return GetLedgerResponse{}, fmt.Errorf("%w: %w", ErrNotFound, err)
	case err != nil:
		return GetLedgerResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return GetLedgerResponse{ledgerFromDomain(ledger)}, nil
}
