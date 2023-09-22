package postgres

import (
	"database/sql"
	"errors"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrorCodeForeignKeyFails = "23503"
	ErrorCodeUniqueFails     = "23505"
)

func ConvertError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNotFound
	}

	psqlErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}

	switch psqlErr.Code {
	case ErrorCodeForeignKeyFails:
		return storage.ErrForeignKeyConstraintFail
	case ErrorCodeUniqueFails:
		return storage.ErrCodeUniqueFails
	default:
		return err
	}
}
