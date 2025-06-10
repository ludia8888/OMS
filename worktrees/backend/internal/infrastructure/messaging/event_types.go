package messaging

import (
	"context"
	"time"

	"github.com/openfoundry/oms/internal/domain/event"
)

// Event represents a domain event
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	EntityID  string                 `json:"entityId"`
	Actor     string                 `json:"actor"`
	Timestamp time.Time              `json:"timestamp"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error
}

// Convert converts internal Event to domain Event
func (e Event) ToDomainEvent() event.Event {
	return event.Event{
		ID:            e.ID,
		EventType:     e.Type,
		AggregateID:   e.EntityID,
		AggregateType: getAggregateType(e.Type),
		Version:       1, // This should be retrieved from the entity
		Timestamp:     e.Timestamp,
		UserID:        e.Actor,
		Data:          e.Data,
		Metadata:      convertMetadata(e.Metadata),
	}
}

func getAggregateType(eventType string) string {
	switch {
	case eventType == EventObjectTypeCreated || eventType == EventObjectTypeUpdated || eventType == EventObjectTypeDeleted:
		return "object_type"
	case eventType == EventLinkTypeCreated || eventType == EventLinkTypeUpdated || eventType == EventLinkTypeDeleted:
		return "link_type"
	default:
		return "unknown"
	}
}

func convertMetadata(metadata map[string]interface{}) map[string]string {
	if metadata == nil {
		return nil
	}
	
	result := make(map[string]string, len(metadata))
	for k, v := range metadata {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}

// Event types
const (
	// Object type events
	EventObjectTypeCreated = "object_type.created"
	EventObjectTypeUpdated = "object_type.updated"
	EventObjectTypeDeleted = "object_type.deleted"

	// Link type events
	EventLinkTypeCreated = "link_type.created"
	EventLinkTypeUpdated = "link_type.updated"
	EventLinkTypeDeleted = "link_type.deleted"

	// Property events
	EventPropertyAdded    = "property.added"
	EventPropertyUpdated  = "property.updated"
	EventPropertyRemoved  = "property.removed"

	// Instance events
	EventInstanceCreated = "instance.created"
	EventInstanceUpdated = "instance.updated"
	EventInstanceDeleted = "instance.deleted"

	// Link events
	EventLinkCreated = "link.created"
	EventLinkUpdated = "link.updated"
	EventLinkDeleted = "link.deleted"
)