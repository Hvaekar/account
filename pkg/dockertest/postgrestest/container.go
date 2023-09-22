package postgrestest

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/dockertest"
	_ "github.com/jackc/pgx/v5/stdlib"
	stddockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

const (
	DefaultStorageUser     = "test_user"
	DefaultStoragePassword = "test_password"
	DefaultStorageDB       = "test_database"
	DefaultSSlMode         = "disable"
	portID                 = "5432/tcp"

	DefaultContainerName = "test-postgres"
)

type Container struct {
	*dockertest.BasicContainer
	user     string
	password string
	db       string
	sslmode  string
}

func NewContainer(opts *stddockertest.RunOptions, user, password, db, sslmode string) *Container {
	opts.Env = append(
		opts.Env,
		fmt.Sprintf("POSTGRES_USER=%s", user),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		fmt.Sprintf("POSTGRES_DB=%s", db),
		fmt.Sprintf("POSTGRES_SSLMODE=%s", sslmode),
	)

	return &Container{
		BasicContainer: dockertest.NewBasicContainer(opts),
		user:           user,
		password:       password,
		db:             db,
		sslmode:        sslmode,
	}
}

func NewDefaultContainer(network *docker.Network) *Container {
	return NewContainer(DefaultRunOptions(network), DefaultStorageUser, DefaultStoragePassword, DefaultStorageDB, DefaultSSlMode)
}

func DefaultRunOptions(network *docker.Network) *stddockertest.RunOptions {
	return dockertest.PrepareRunOptions(
		&stddockertest.RunOptions{
			Name:       DefaultContainerName,
			Repository: "postgres",
			Tag:        "latest",
			NetworkID:  network.ID,
		},
		portID,
	)
}

func (c *Container) WaitReady(ctx context.Context) error {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			db, err := sql.Open("pgx", c.Dsn())
			if err != nil {
				timer.Reset(time.Second)
				continue
			}

			if err := db.Ping(); err != nil {
				timer.Reset(time.Second)
				continue
			}

			return db.Close()
		}
	}
}

func (c *Container) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.GetUser(),
		c.GetPassword(),
		c.GetHost(),
		c.GetPort(),
		c.GetDB(),
		c.GetSSLMode(),
	)
}

func (c *Container) GetHost() string {
	return dockertest.Host(c.Resource)
}

func (c *Container) GetUser() string {
	return c.user
}

func (c *Container) GetPassword() string {
	return c.password
}

func (c *Container) GetPort() string {
	return c.Resource.GetPort(portID)
}

func (c *Container) GetDB() string {
	return c.db
}

func (c *Container) GetSSLMode() string {
	return c.sslmode
}
