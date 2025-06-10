package event

import (
	"context"
	"time"
)

// Event represents a domain event
type Event struct {
	ID            string      `json:"id"`
	EventType     string      `json:"eventType"`
	AggregateID   string      `json:"aggregateId"`
	AggregateType string      `json:"aggregateType"`
	Version       int         `json:"version"`
	Timestamp     time.Time   `json:"timestamp"`
	UserID        string      `json:"userId"`
	Data          interface{} `json:"data"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error
}

// EventStore defines the interface for storing events
type EventStore interface {
	Save(ctx context.Context, event Event) error
	GetByAggregateID(ctx context.Context, aggregateID string) ([]Event, error)
	GetByEventType(ctx context.Context, eventType string, limit int) ([]Event, error)
}