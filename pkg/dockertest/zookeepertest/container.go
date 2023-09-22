package zookeepertest

import (
	"context"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/dockertest"
	"github.com/go-zookeeper/zk"
	stddockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

const (
	DefaultContainerName = "test-zookeeper"

	DefaultZooMyID   = "1"
	DefaultZooServes = "server.1=" + DefaultContainerName + ":2888:3888;" + DefaultPort
	DefaultPort      = "2181"
	portID           = "2181/tcp"
)

type Container struct {
	*dockertest.BasicContainer
	id     string
	serves string
}

func NewContainer(opts *stddockertest.RunOptions, id, port, serves string) *Container {
	opts.Env = append(
		opts.Env,
		fmt.Sprintf("ZOO_MY_ID=%s", id),
		fmt.Sprintf("ZOO_PORT=%s", port),
		fmt.Sprintf("ZOO_SERVERS=%s", serves),
	)

	return &Container{
		BasicContainer: dockertest.NewBasicContainer(opts),
		id:             id,
		serves:         serves,
	}
}

func NewDefaultContainer(network *docker.Network) *Container {
	return NewContainer(DefaultRunOptions(network), DefaultZooMyID, DefaultPort, DefaultZooServes)
}

func DefaultRunOptions(network *docker.Network) *stddockertest.RunOptions {
	return dockertest.PrepareRunOptions(
		&stddockertest.RunOptions{
			Name:       DefaultContainerName,
			Repository: "zookeeper",
			Tag:        "latest",
			Hostname:   DefaultContainerName,
			NetworkID:  network.ID,
			ExposedPorts: []string{
				portID,
			},
		},
	)
}

func (c *Container) WaitReady(ctx context.Context) error {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", c.GetHost(), c.GetPort())}, 10*time.Second)
			if err != nil {
				timer.Reset(time.Second)
				continue
			}
			defer conn.Close()

			return nil
		}
	}
}

func (c *Container) GetHost() string {
	return dockertest.Host(c.Resource)
}

func (c *Container) GetID() string {
	return c.id
}

func (c *Container) GetPort() string {
	return c.Resource.GetPort(portID)
}

func (c *Container) GetServes() string {
	return c.serves
}
