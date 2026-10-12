package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type accountRepository struct {
	tx pgx.Tx
}

func (r accountRepository) Insert(ctx context.Context, a domain.Account) error {
	const query = "INSERT INTO accounts (account_id, ledger_id, type, code, name, created_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.tx.Exec(
		ctx,
		query,
		a.AccountID,
		a.LedgerID,
		a.Type,
		a.Code,
		a.Name,
		a.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "accounts_ledger_id_code_unique" {
				return fmt.Errorf("%w: code: %q, ledger id: %s", domain.ErrDuplicateAccountCodeInLedger, a.Code, a.LedgerID)
			}
		}
		// Default error message.
		return fmt.Errorf("insert account: %w", err)
	}

	return nil
}
