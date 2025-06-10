package messaging

import (
	"context"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	EventObjectTypeCreated EventType = "ObjectTypeCreated"
	EventObjectTypeUpdated EventType = "ObjectTypeUpdated"
	EventObjectTypeDeleted EventType = "ObjectTypeDeleted"
	EventLinkTypeCreated   EventType = "LinkTypeCreated"
	EventLinkTypeUpdated   EventType = "LinkTypeUpdated"
	EventLinkTypeDeleted   EventType = "LinkTypeDeleted"
)

// Event represents a domain event
type Event struct {
	ID            string                 `json:"id"`
	Type          EventType              `json:"type"`
	EntityID      string                 `json:"entityId"`
	Actor         string                 `json:"actor"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          interface{}            `json:"data"`
	Metadata      map[string]interface{} `json:"metadata"`
	CorrelationID string                 `json:"correlationId,omitempty"`
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	// Publish publishes an event
	Publish(ctx context.Context, event Event) error
	
	// PublishBatch publishes multiple events
	PublishBatch(ctx context.Context, events []Event) error
	
	// Close closes the publisher connection
	Close() error
}

// EventSubscriber defines the interface for subscribing to events
type EventSubscriber interface {
	// Subscribe subscribes to events
	Subscribe(ctx context.Context, eventTypes []EventType, handler EventHandler) error
	
	// Unsubscribe unsubscribes from events
	Unsubscribe() error
	
	// Close closes the subscriber connection
	Close() error
}

// EventHandler is a function that handles events
type EventHandler func(ctx context.Context, event Event) error