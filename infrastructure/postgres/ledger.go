package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/banovic/saldo/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

// ledgerRepository implements ledger related db operations.
type ledgerRepository struct {
	tx pgx.Tx
}

// ledgerRow represents row of data from database, before it is converted into domain object.
type ledgerRow struct {
	LedgerID           uuid.UUID
	Name               domain.LedgerName
	FunctionalCurrency domain.Currency
	ReportingTimeZone  domain.TimeZone
	CreatedAt          time.Time
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

	// Handle error.
	if err != nil {
		// Error is postgres error, with all relevant details.
		if pgErr, ok := pgError(err); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				// There can be multiple unique constraints, find out which.
				if pgErr.ConstraintName == "ledgers_name_unique" {
					return fmt.Errorf("%w: %q", domain.ErrDuplicateLedgerName, l.Name)
				}
			}
		}
		return fmt.Errorf("insert ledger: %w", err)
	}
	return nil
}

// Get one ledger by its id (which is primary key), returns error on failures.
func (lr ledgerRepository) Get(ctx context.Context, id domain.LedgerID) (domain.Ledger, error) {
	const query = "SELECT ledger_id, name, functional_currency, reporting_time_zone, created_at FROM ledgers WHERE ledger_id = $1"
	row, err := lr.tx.Query(ctx, query, id)
	if err != nil {
		return domain.Ledger{}, fmt.Errorf("ledger get: %w", err)
	}
	r, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[ledgerRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Ledger{}, fmt.Errorf("%w: %v", domain.ErrLedgerNotFound, id)
	}
	if err != nil {
		return domain.Ledger{}, fmt.Errorf("ledger get: %w", err)
	}
	return r.toDomain(), nil
}

// List all ledgers, returns error on failures.
func (lr ledgerRepository) List(ctx context.Context) ([]domain.Ledger, error) {
	const query = "SELECT ledger_id, name, functional_currency, reporting_time_zone, created_at FROM ledgers"
	rows, err := lr.tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ledgers list: %w", err)
	}
	rs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ledgerRow])
	if err != nil {
		return nil, fmt.Errorf("list ledgers: %w", err)
	}
	drs := make([]domain.Ledger, 0, len(rs))
	for _, r := range rs {
		drs = append(drs, r.toDomain())
	}
	return drs, nil
}

// toDomain converts row returned from database, into domain object.
func (r ledgerRow) toDomain() domain.Ledger {
	return domain.Ledger{
		LedgerID:           domain.LedgerID{UUID: r.LedgerID},
		Name:               r.Name,
		FunctionalCurrency: r.FunctionalCurrency,
		ReportingTimeZone:  r.ReportingTimeZone,
		CreatedAt:          r.CreatedAt,
	}
}
