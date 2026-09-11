package postgres

import (
	"context"
	"errors"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgx/v5"
)

type ledgerRepository struct {
	tx pgx.Tx
}

func (lr ledgerRepository) Insert(ctx context.Context, l domain.Ledger) error {
	return errors.New("TODO")
}
