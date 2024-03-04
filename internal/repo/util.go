package repo

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lmnq/test-thai/internal/errs"
)

func isUniqueConstraintError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgErr.Code == errs.UniqueConstraintCode
}
