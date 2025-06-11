package eventstore

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// MySQLEventStore implements Store backed by a MySQL table named `events`.
type MySQLEventStore struct {
	db *gorm.DB
}

// NewMySQL creates a new MySQLEventStore.
func NewMySQL(db *gorm.DB) *MySQLEventStore {
	return &MySQLEventStore{db: db}
}

// Append stores the event type and JSON payload in the events table.
type EventRecord struct {
	ID        int `gorm:"primaryKey"`
	Type      string
	Payload   string
	CreatedAt time.Time
}

func (EventRecord) TableName() string {
	return "events"
}

func (s *MySQLEventStore) Append(eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	record := EventRecord{
		Type:      eventType,
		Payload:   string(data),
		CreatedAt: time.Now(),
	}
	return s.db.Create(&record).Error
}
