package broker

type MessageBroker interface {
	AddConsumer() error
	Consume() error
	AddProducer() error
	SendMessage(topic string, message any) error
	GetTopic(name string) string
}
