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
		return ListLedgersResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	dtoLedgers := make([]Ledger, 0, len(ledgers))
	for _, ledger := range ledgers {
		dtoLedgers = append(dtoLedgers, ledgerFromDomain(ledger))
	}
	return ListLedgersResponse{Ledgers: dtoLedgers}, nil
}
