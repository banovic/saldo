package app

import (
	"context"
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

	if err != nil {
		return GetLedgerResponse{}, fmt.Errorf("get ledger: %w", err)
	}

	return GetLedgerResponse{
		LedgerID:           ledger.LedgerID,
		Name:               ledger.Name,
		FunctionalCurrency: ledger.FunctionalCurrency,
		ReportingTimeZone:  ledger.ReportingTimeZone,
		CreatedAt:          ledger.CreatedAt,
	}, nil
}
