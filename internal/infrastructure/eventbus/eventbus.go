package eventbus

import (
	"sync"

	"CardGameDB/internal/infrastructure/eventstore"
)

// EventBus is a simple event dispatcher

type EventBus struct {
	listeners map[string][]func(interface{})
	store     eventstore.Store
	mu        sync.RWMutex
}

// New creates EventBus
func New(store eventstore.Store) *EventBus {
	return &EventBus{
		listeners: make(map[string][]func(interface{})),
		store:     store,
	}
}

// Subscribe to event by name
func (eb *EventBus) Subscribe(event string, handler func(interface{})) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners[event] = append(eb.listeners[event], handler)
}

// Publish event by name
func (eb *EventBus) Publish(event string, payload interface{}) {
	if eb.store != nil {
		// ignore store errors for now
		_ = eb.store.Append(event, payload)
	}
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if handlers, ok := eb.listeners[event]; ok {
		for _, h := range handlers {
			go h(payload)
		}
	}
}
