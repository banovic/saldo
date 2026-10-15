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
	"github.com/jackc/pgx/v5/pgconn"
)

type accountRepository struct {
	tx pgx.Tx
}

type accountRow struct {
	AccountID uuid.UUID
	LedgerID  uuid.UUID
	Type      domain.AccountType
	Code      domain.AccountCode
	Name      domain.AccountName
	CreatedAt time.Time
}

func (row accountRow) toDomain() domain.Account {
	return domain.Account{
		AccountID: domain.AccountID{UUID: row.AccountID},
		LedgerID:  domain.LedgerID{UUID: row.LedgerID},
		Type:      row.Type,
		Code:      row.Code,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
	}
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

func (r accountRepository) Get(ctx context.Context, id domain.AccountID) (domain.Account, error) {
	const query = "SELECT account_id, ledger_id, type, code, name, created_at FROM accounts WHERE account_id = $1"
	row, err := r.tx.Query(ctx, query, id)
	if err != nil {
		return domain.Account{}, fmt.Errorf("account get: %w", err)
	}
	res, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[accountRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Account{}, fmt.Errorf("%w: %v", domain.ErrAccountNotFound, id)
	}
	if err != nil {
		return domain.Account{}, fmt.Errorf("account get: %w", err)
	}
	return res.toDomain(), nil
}

func (r accountRepository) List(ctx context.Context, ledgerID domain.LedgerID) ([]domain.Account, error) {
	const query = "SELECT account_id, ledger_id, type, code, name, created_at FROM accounts WHERE ledger_id = $1 ORDER BY created_at DESC"
	rows, err := r.tx.Query(ctx, query, ledgerID)
	if err != nil {
		return nil, fmt.Errorf("account list: %w", err)
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[accountRow])
	if err != nil {
		return nil, fmt.Errorf("account list: %w", err)
	}

	d := make([]domain.Account, 0, len(res))
	for _, o := range res {
		d = append(d, o.toDomain())
	}
	return d, nil
}
