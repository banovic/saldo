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

func (s *Service) ListLedgers(ctx context.Context, req ListLedgersRequest) (ListLedgersResponse, error) {
	// Get data from db.
	ledgers, err := executeValue(ctx, s.uow, func(r Repositories) ([]domain.Ledger, error) {
		return r.Ledger.List(ctx)
	})

	// Map err.
	if err != nil {
		return ListLedgersResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Map result.
	return ListLedgersResponse{ledgersFromDomain(ledgers)}, nil
}
