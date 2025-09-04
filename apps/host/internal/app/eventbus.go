package app

import "sync"

// EventBus is a simple pub-sub event bus
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan any
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan any),
	}
}

// Subscribe subscribes to an event and returns a channel to receive events
func (b *EventBus) Subscribe(topic string) <-chan any {
	ch := make(chan any, 1)
	b.mu.Lock()

	b.subscribers[topic] = append(b.subscribers[topic], ch)

	b.mu.Unlock()

	return ch
}

// Publish publishes an event to the bus
func (b *EventBus) Publish(topic string, data any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers[topic] {
		select {
		case ch <- data:
		default:
		}
	}
}
