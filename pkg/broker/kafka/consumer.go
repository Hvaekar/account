package kafka

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/IBM/sarama"
)

func (b *MessageBroker) AddConsumer() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(b.cfg.Brokers, config)
	if err != nil {
		return err
	}

	b.consumer = consumer

	return nil
}

func (b *MessageBroker) Consume() error {
	defer func() {
		if err := b.consumer.Close(); err != nil {
			b.log.Errorf("close consumer: %s", err.Error())
		}
	}()

	accountCreateTopic := b.GetTopic(broker.AccountCreateKey)
	accountCreate, err := b.consumer.ConsumePartition(accountCreateTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer func() {
		if err := accountCreate.Close(); err != nil {
			b.log.Errorf("close account create partition: %s", err.Error())
		}
	}()

	for {
		select {
		case <-accountCreate.Messages():
			//b.log.Info("received account create event") // do something in this microservice after get message
		}
	}
}
