package postgres

import (
	"context"
	"errors"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgx/v5"
)

type accountRepository struct {
	tx pgx.Tx
}

func (ar accountRepository) Insert(ctx context.Context, a domain.Account) error {
	return errors.New("TODO")
}
