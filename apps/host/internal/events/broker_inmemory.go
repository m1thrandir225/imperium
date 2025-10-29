package events

import "sync"

// InMemoryBroker is a simple pub-sub event broker
type InMemoryBroker struct {
	mu          sync.RWMutex
	subscribers map[string][]chan any
}

// NewInMemoryBroker creates a new event bus
func NewInMemoryBroker() EventBroker {
	return &InMemoryBroker{
		subscribers: make(map[string][]chan any),
	}
}

// Subscribe subscribes to an event and returns a channel to receive events
func (b *InMemoryBroker) Subscribe(topic string) <-chan any {
	ch := make(chan any, 1)
	b.mu.Lock()

	b.subscribers[topic] = append(b.subscribers[topic], ch)

	b.mu.Unlock()

	return ch
}

// Publish publishes an event to the bus
func (b *InMemoryBroker) Publish(topic string, data any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers[topic] {
		select {
		case ch <- data:
		default:
		}
	}
}
