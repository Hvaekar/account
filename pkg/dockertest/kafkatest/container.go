package kafkatest

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/dockertest"
	"github.com/IBM/sarama"
	stddockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

const (
	DefaultContainerName = "test-kafka"
	portID               = "29092/tcp"
)

type Container struct {
	*dockertest.BasicContainer
}

func NewContainer(opts *stddockertest.RunOptions) *Container {
	return &Container{
		BasicContainer: dockertest.NewBasicContainer(opts),
	}
}

func NewDefaultContainer(network *docker.Network) *Container {
	return NewContainer(DefaultRunOptions(network))
}

func DefaultRunOptions(network *docker.Network) *stddockertest.RunOptions {
	return &stddockertest.RunOptions{
		Name:       DefaultContainerName,
		Repository: "wurstmeister/kafka",
		Tag:        "latest",
		NetworkID:  network.ID,
		Hostname:   DefaultContainerName,
		Env: []string{
			"KAFKA_CREATE_TOPICS=domain.test:1:1:compact",
			"KAFKA_ADVERTISED_LISTENERS=INSIDE://" + DefaultContainerName + ":9092,OUTSIDE://localhost:29092",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT",
			"KAFKA_LISTENERS=INSIDE://0.0.0.0:9092,OUTSIDE://0.0.0.0:29092",
			"KAFKA_ZOOKEEPER_CONNECT=test-zookeeper:2181",
			"KAFKA_INTER_BROKER_LISTENER_NAME=INSIDE",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			portID: {{HostIP: "localhost", HostPort: portID}},
		},
		ExposedPorts: []string{portID},
	}
}

func (c *Container) WaitReady(ctx context.Context) error {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			config := sarama.NewConfig()
			config.Producer.Return.Successes = true

			producer, err := sarama.NewAsyncProducer(c.GetBrokers(), config)
			if err != nil {
				timer.Reset(time.Second)
				continue
			}
			defer producer.Close()

			message := &sarama.ProducerMessage{
				Topic: "domain.test",
				Value: sarama.StringEncoder("Hello World"),
			}

			producer.Input() <- message

			select {
			case <-producer.Successes():
				return nil
			case <-producer.Errors():
				timer.Reset(time.Second)
				continue
			}
		}
	}
}

func (c *Container) GetBrokers() []string {
	return []string{c.GetHost() + ":" + c.GetPort()}
}

func (c *Container) GetHost() string {
	return dockertest.Host(c.Resource)
}

func (c *Container) GetPort() string {
	return c.Resource.GetPort(portID)
}
