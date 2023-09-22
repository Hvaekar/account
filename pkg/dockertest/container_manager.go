package dockertest

import (
	"context"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"golang.org/x/sync/errgroup"
)

type ContainerManager struct {
	containers []Container
	started    bool
	pool       *dockertest.Pool
	network    *docker.Network
	log        logger.Logger
}

func NewContainerManager(log logger.Logger) *ContainerManager {
	return &ContainerManager{containers: make([]Container, 0), log: log}
}

func (m *ContainerManager) CreatePool() error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("dockertest new pool: %w", err)
	}

	m.pool = pool

	return nil
}

func (m *ContainerManager) AddNetwork(name string) error {
	network, err := m.pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: name})
	if err != nil {
		return fmt.Errorf("create network: %w", err)
	}

	m.network = network

	return nil
}

func (m *ContainerManager) RunAndWaitReady(ctx context.Context) error {
	if m.started {
		return nil
	}

	if err := m.run(ctx); err != nil {
		return err
	}

	if err := m.waitReady(ctx); err != nil {
		return err
	}

	m.started = true

	return nil
}

func (m *ContainerManager) run(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, container := range m.containers {
		m.log.Infof("run container '%s'", container.GetName())

		container := container

		g.Go(func() error {
			if err := container.Run(gCtx, m.pool); err != nil {
				return err
			}

			return nil
		})
	}

	return g.Wait()
}

func (m *ContainerManager) waitReady(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, container := range m.containers {
		m.log.Infof("checking container '%s'", container.GetName())

		container := container

		g.Go(func() error {
			if err := container.WaitReady(gCtx); err != nil {
				return err
			}

			m.log.Infof("'%s' checked", container.GetName())

			return nil
		})
	}

	return g.Wait()
}

func (m *ContainerManager) AddContainer(container Container) {
	m.containers = append(m.containers, container)
}

func (m *ContainerManager) Stop(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, container := range m.containers {
		m.log.Infof("stopping container '%s'", container.GetName())

		container := container

		g.Go(func() error {
			if err := container.Stop(gCtx); err != nil {
				return err
			}

			return nil
		})
	}

	m.started = false

	return g.Wait()
}

func (m *ContainerManager) GetNetwork() *docker.Network {
	return m.network
}

func (m *ContainerManager) RemoveNetwork(id string) error {
	if err := m.pool.Client.RemoveNetwork(id); err != nil {
		return fmt.Errorf("remove network: %w", err)
	}

	return nil
}
