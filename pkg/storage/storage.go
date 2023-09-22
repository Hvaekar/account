package storage

import (
	"context"
	"database/sql"
)

type Storage interface {
	Connect() error
	Migrate() error
	Close() error
	GetDB() *sql.DB
	TruncateTables(c context.Context, tables ...string) error
}
