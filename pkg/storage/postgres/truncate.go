package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Postgres) TruncateTables(c context.Context, tables ...string) error {
	for i := range tables {
		if _, err := s.DB.ExecContext(c, "TRUNCATE TABLE "+tables[i]+" CASCADE"); err != nil {
			return fmt.Errorf("truncating table %s: %w", tables[i], err)
		}
	}

	return nil
}
