package storage

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrForeignKeyConstraintFail = errors.New("foreign key constraint fail")
	ErrCodeUniqueFails          = errors.New("unique constraint fail")
	ErrIncorrectPassword        = errors.New("incorrect password")
)
