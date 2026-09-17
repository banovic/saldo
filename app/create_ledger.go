package app

import (
	"context"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type CreateLedgerRequest struct {
	Name               string
	FunctionalCurrency string
	ReportingTimeZone  string
}

type CreateLedgerResponse struct {
	LedgerID domain.LedgerID
}

func (svc *Service) CreateLedger(ctx context.Context, req CreateLedgerRequest) (CreateLedgerResponse, error) {
	ledger := domain.Ledger{
		LedgerID:           domain.LedgerID{UUID: svc.generateUUID()},
		Name:               domain.LedgerName(req.Name),
		FunctionalCurrency: domain.Currency(req.FunctionalCurrency),
		ReportingTimeZone:  domain.TimeZone(req.ReportingTimeZone),
		CreatedAt:          svc.utcTimestamp(),
	}

	if err := ledger.Validate(); err != nil {
		return CreateLedgerResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	if err := svc.uow.Execute(ctx, func(r Repositories) error {
		return r.Ledger.Insert(ctx, ledger)
	}); err != nil {
		return CreateLedgerResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return CreateLedgerResponse{LedgerID: ledger.LedgerID}, nil
}
