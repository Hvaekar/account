package dockertest

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
)

type Container interface {
	Run(ctx context.Context, pool *dockertest.Pool) error
	Stop(ctx context.Context) error
	WaitReady(ctx context.Context) error
	GetName() string
	GetResource() *dockertest.Resource
}

type BasicContainer struct {
	Options  *dockertest.RunOptions
	Resource *dockertest.Resource
}

func NewBasicContainer(options *dockertest.RunOptions) *BasicContainer {
	return &BasicContainer{Options: options}
}

func (c *BasicContainer) Run(_ context.Context, pool *dockertest.Pool) error {
	if resource, ok := pool.ContainerByName(c.Options.Name); ok {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("pool purge: %w", err)
		}
	}

	r, err := pool.RunWithOptions(c.Options)
	if err != nil {
		return fmt.Errorf("run docker container: %w", err)
	}

	c.Resource = r

	return nil
}

func (c *BasicContainer) Stop(_ context.Context) error {
	return c.Resource.Close()
}

func (c *BasicContainer) GetName() string {
	return c.Options.Name
}

func (c *BasicContainer) GetResource() *dockertest.Resource {
	return c.Resource
}
