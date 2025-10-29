// Package events provides a definition for an event pub-sub system
// either in-memory or other implementations used for the Event Driven
// Architecture of the UI Layer.
package events

type EventBroker interface {
	Subscribe(topic string) <-chan any
	Publish(topic string, data any)
}
