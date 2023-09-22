package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

func (b *MessageBroker) AddProducer() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(b.cfg.Brokers, config)
	if err != nil {
		return err
	}

	b.producer = producer

	return nil
}

func (b *MessageBroker) SendMessage(topic string, message any) error {
	topic = b.GetTopic(topic)
	if topic == "" {
		return fmt.Errorf("empty topic name")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(body),
	}

	b.producer.Input() <- msg

	select {
	case <-b.producer.Successes():
		return nil
	case err := <-b.producer.Errors():
		return err.Err
	}
}
