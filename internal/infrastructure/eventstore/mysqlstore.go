package eventstore

import (
	"database/sql"
	"encoding/json"
	"time"
)

// MySQLEventStore implements Store backed by a MySQL table named `events`.
type MySQLEventStore struct {
	db *sql.DB
}

// NewMySQL creates a new MySQLEventStore.
func NewMySQL(db *sql.DB) *MySQLEventStore {
	return &MySQLEventStore{db: db}
}

// Append stores the event type and JSON payload in the events table.
func (s *MySQLEventStore) Append(eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO events (type, payload, created_at) VALUES (?, ?, ?)", eventType, string(data), time.Now())
	return err
}
