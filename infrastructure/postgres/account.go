package postgres

import (
	"errors"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgx/v5"
)

type accountRepository struct {
	db *pgx.Tx
}

func (lr accountRepository) Insert(l domain.Account) error {
	return errors.New("TODO")
}
