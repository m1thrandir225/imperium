package event_broker

type EventBroker interface {
	Subscribe(topic string) <-chan any
	Publish(topic string, data any)
}
