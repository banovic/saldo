package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func pgError(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr, true
	}
	return nil, false
}
