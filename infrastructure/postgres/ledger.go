package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgerrcode"
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
		if pgErr.Code == pgerrcode.UniqueViolation {
			// There can be multiple unique constraints, find out which.
			if pgErr.ConstraintName == "ledgers_name_unique" {
				return fmt.Errorf("%w: %q", domain.ErrDuplicateLedgerName, l.Name)
			}
		}
		return fmt.Errorf("insert ledger error")
	}
	return nil
}

func (lr ledgerRepository) Get(ctx context.Context, id domain.LedgerID) (domain.Ledger, error) {
	const query = "SELECT * FROM ledgers WHERE ledger_id = $1"
	row, err := lr.tx.Query(ctx, query, id)
	if err != nil {
		return domain.Ledger{}, nil
	}
	ledger, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[domain.Ledger])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Ledger{}, fmt.Errorf("%w: %v", domain.ErrLedgerNotFound, id)
	}
	if err != nil {
		return domain.Ledger{}, fmt.Errorf("ledger get: %w", err)
	}
	return ledger, nil
}
