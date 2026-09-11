package postgres

import (
	"errors"

	"github.com/banovic/saldo/domain"
)

type ledgerRepository struct {
	// db ???
}

func (lr ledgerRepository) Insert(l domain.Ledger) error {
	return errors.New("TODO")
}
