package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Hvaekar/med-account/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"time"
)

type Postgres struct {
	DB  *sql.DB
	cfg *config.Postgres
}

func NewPostgres(cfg *config.Postgres) *Postgres {
	return &Postgres{cfg: cfg}
}

func (s *Postgres) Connect() error {
	db, err := sql.Open(s.cfg.Driver, s.dsn())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(s.cfg.MaxOpenConns)
	db.SetConnMaxLifetime(s.cfg.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(s.cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(s.cfg.ConnMaxIdleTime * time.Second)

	if err := db.Ping(); err != nil {
		return err
	}

	s.DB = db

	return nil
}

func (s *Postgres) GetDB() *sql.DB {
	return s.DB
}

func (s *Postgres) Migrate() error {
	m, err := migrate.New("file://migrations", s.dsn())
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *Postgres) Close() error {
	return s.DB.Close()
}

func (s *Postgres) dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		s.cfg.User,
		s.cfg.Password,
		s.cfg.Host,
		s.cfg.Port,
		s.cfg.DB,
		s.cfg.SSLMode,
	)
}
