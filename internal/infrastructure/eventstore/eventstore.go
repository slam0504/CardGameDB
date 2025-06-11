package eventstore

// Store defines an append-only event store.
type Store interface {
	Append(eventType string, payload interface{}) error
}
