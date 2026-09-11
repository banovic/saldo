package app

import (
	"context"

	"github.com/banovic/saldo/domain"
)

type CreateLedgerRequest struct {
	LedgerID           domain.LedgerID
	Name               string
	FunctionalCurrency domain.Currency
}

type CreateLedgerResponse struct {
	LedgerID domain.LedgerID
}

func (*Service) CreateLedger(ctx context.Context, req CreateLedgerRequest) (CreateLedgerResponse, error) {
	return CreateLedgerResponse{}, nil
}
