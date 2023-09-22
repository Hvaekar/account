package kafka

import (
	"github.com/Hvaekar/med-account/config"
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/IBM/sarama"
)

type MessageBroker struct {
	producer sarama.AsyncProducer
	consumer sarama.Consumer
	cfg      *config.Kafka
	log      logger.Logger
}

func NewMessageBroker(cfg *config.Kafka, log logger.Logger) broker.MessageBroker {
	return &MessageBroker{
		cfg: cfg,
		log: log,
	}
}
