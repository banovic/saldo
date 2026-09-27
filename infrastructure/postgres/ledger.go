package postgres

import (
	"context"
	"fmt"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgx/v5"
)

// ledgerRepository implements ledger related db operations.
type ledgerRepository struct {
	tx pgx.Tx
}

// Insert inserts new Ledger into database, and returns error on failure.
func (lr ledgerRepository) Insert(ctx context.Context, l domain.Ledger) error {
	const query = "INSERT INTO ledgers (ledger_id, name, functional_currency, reporting_time_zone, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := lr.tx.Exec(
		ctx,
		query,
		l.LedgerID,
		l.Name,
		l.FunctionalCurrency,
		l.ReportingTimeZone,
		l.CreatedAt,
	)

	// Check what can go wrong with postgres - schema violation.
	if pgErr, ok := pgError(err); ok {
		if pgErr.Code == "23505" {
			// There can be multiple unique constraints, find out which.
			if pgErr.ConstraintName == "ledgers_name_unique" {
				return fmt.Errorf("%w: %q", domain.ErrDuplicateLedgerName, l.Name)
			}
		}
		return fmt.Errorf("insert ledger error")
	}
	return nil
}
