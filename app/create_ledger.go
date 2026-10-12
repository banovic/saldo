package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type CreateLedgerRequest struct {
	Name               string
	FunctionalCurrency string
	ReportingTimeZone  string
}

type CreateLedgerResponse struct{ Ledger }

func (s *Service) CreateLedger(ctx context.Context, req CreateLedgerRequest) (CreateLedgerResponse, error) {
	ledger := domain.Ledger{
		LedgerID:           domain.LedgerID{UUID: s.generateUUID()},
		Name:               domain.LedgerName(req.Name),
		FunctionalCurrency: domain.Currency(req.FunctionalCurrency),
		ReportingTimeZone:  domain.TimeZone(req.ReportingTimeZone),
		CreatedAt:          s.utcTimestamp(),
	}

	if err := ledger.Validate(); err != nil {
		return CreateLedgerResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	err := s.uow.Execute(ctx, func(r Repositories) error {
		return r.Ledger.Insert(ctx, ledger)
	})
	switch {
	case errors.Is(err, domain.ErrDuplicateLedgerName):
		return CreateLedgerResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	case err != nil:
		return CreateLedgerResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return CreateLedgerResponse{ledgerFromDomain(ledger)}, nil
}
