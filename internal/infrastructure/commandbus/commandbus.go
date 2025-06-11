package commandbus

import (
	"encoding/json"
	"sync"

	"CardGameDB/internal/infrastructure/eventstore"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/kafka"
)

// CommandBus is a simple event dispatcher

type CommandBus struct {
	listeners map[string][]func(interface{})
	store     eventstore.Store
	publisher message.Publisher
	mu        sync.RWMutex
}

// New creates CommandBus
func New(store eventstore.Store) *CommandBus {
	return &CommandBus{
		listeners: make(map[string][]func(interface{})),
		store:     store,
	}
}

// NewKafka creates CommandBus that also publishes messages to Kafka using Watermill
func NewKafka(store eventstore.Store, brokers []string) (*CommandBus, error) {
	pub, err := kafka.NewPublisher(kafka.PublisherConfig{Brokers: brokers}, watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, err
	}
	return &CommandBus{
		listeners: make(map[string][]func(interface{})),
		store:     store,
		publisher: pub,
	}, nil
}

// Subscribe to event by name
func (eb *CommandBus) Subscribe(event string, handler func(interface{})) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners[event] = append(eb.listeners[event], handler)
}

// Publish event by name
func (eb *CommandBus) Publish(event string, payload interface{}) {
	if eb.store != nil {
		// ignore store errors for now
		_ = eb.store.Append(event, payload)
	}
	if eb.publisher != nil {
		data, _ := json.Marshal(payload)
		msg := message.NewMessage(watermill.NewUUID(), data)
		if err := eb.publisher.Publish(event, msg); err != nil {
			// ignore publish errors
		}
	}
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if handlers, ok := eb.listeners[event]; ok {
		for _, h := range handlers {
			go h(payload)
		}
	}
}
