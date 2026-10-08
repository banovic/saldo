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

func (svc *Service) GetLedger(ctx context.Context, req GetLedgerRequest) (GetLedgerResponse, error) {
	ledgerID, err := domain.ParseLedgerID(req.LedgerID)
	if err != nil {
		return GetLedgerResponse{}, fmt.Errorf("%w: %q: %w", ErrInvalidInput, req.LedgerID, err)
	}

	var ledger domain.Ledger

	err = svc.uow.Execute(ctx, func(r Repositories) error {
		var err error
		ledger, err = r.Ledger.Get(ctx, ledgerID)
		return err
	})

	switch {
	case errors.Is(err, domain.ErrLedgerNotFound):
		return GetLedgerResponse{}, fmt.Errorf("%w: %w", ErrNotFound, err)
	case err != nil:
		return GetLedgerResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return GetLedgerResponse{ledgerFromDomain(ledger)}, nil
}
