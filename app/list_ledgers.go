package app

import (
	"context"
	"fmt"

	"github.com/banovic/saldo/domain"
)

type ListLedgersRequest struct{}

type ListLedgersResponse struct {
	Ledgers []Ledger
}

func (svc *Service) ListLedgers(ctx context.Context, req ListLedgersRequest) (ListLedgersResponse, error) {
	var ledgers []domain.Ledger

	err := svc.uow.Execute(ctx, func(r Repositories) error {
		var err error
		ledgers, err = r.Ledger.List(ctx)
		return err
	})
	if err != nil {
		return ListLedgersResponse{}, fmt.Errorf("list ledgers: %w", err)
	}
	dtoLedgers := make([]Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		dtoLedgers = append(dtoLedgers, Ledger{
			LedgerID:           ledger.LedgerID,
			Name:               ledger.Name,
			FunctionalCurrency: ledger.FunctionalCurrency,
			ReportingTimeZone:  ledger.ReportingTimeZone,
			CreatedAt:          ledger.CreatedAt,
		})
	}
	return ListLedgersResponse{Ledgers: dtoLedgers}, nil
}
