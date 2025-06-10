package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openfoundry/oms/internal/domain/event"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaPublisher implements the EventPublisher interface using Kafka
type KafkaPublisher struct {
	writer *kafka.Writer
	logger *zap.Logger
}

// NewKafkaPublisher creates a new Kafka event publisher
func NewKafkaPublisher(brokers []string, topic string, logger *zap.Logger) *KafkaPublisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Compression:  kafka.Snappy,
		Async:        false, // Synchronous for reliability
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  3,
		Logger:       kafka.LoggerFunc(logger.Sugar().Debugf),
		ErrorLogger:  kafka.LoggerFunc(logger.Sugar().Errorf),
	}

	return &KafkaPublisher{
		writer: writer,
		logger: logger,
	}
}

// Publish publishes an event to Kafka
func (p *KafkaPublisher) Publish(ctx context.Context, evt event.Event) error {
	// Marshal event data
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(evt.AggregateID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(evt.EventType)},
			{Key: "aggregate_type", Value: []byte(evt.AggregateType)},
			{Key: "version", Value: []byte(fmt.Sprintf("%d", evt.Version))},
		},
		Time: evt.Timestamp,
	}

	// Publish to Kafka
	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("event_id", evt.ID),
			zap.String("event_type", evt.EventType),
			zap.String("aggregate_id", evt.AggregateID),
			zap.Error(err))
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("Event published",
		zap.String("event_id", evt.ID),
		zap.String("event_type", evt.EventType),
		zap.String("aggregate_id", evt.AggregateID))

	return nil
}

// PublishBatch publishes multiple events to Kafka
func (p *KafkaPublisher) PublishBatch(ctx context.Context, events []event.Event) error {
	messages := make([]kafka.Message, 0, len(events))

	for _, evt := range events {
		data, err := json.Marshal(evt)
		if err != nil {
			return fmt.Errorf("failed to marshal event %s: %w", evt.ID, err)
		}

		message := kafka.Message{
			Key:   []byte(evt.AggregateID),
			Value: data,
			Headers: []kafka.Header{
				{Key: "event_type", Value: []byte(evt.EventType)},
				{Key: "aggregate_type", Value: []byte(evt.AggregateType)},
				{Key: "version", Value: []byte(fmt.Sprintf("%d", evt.Version))},
			},
			Time: evt.Timestamp,
		}

		messages = append(messages, message)
	}

	// Publish batch to Kafka
	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		p.logger.Error("Failed to publish event batch",
			zap.Int("batch_size", len(events)),
			zap.Error(err))
		return fmt.Errorf("failed to publish event batch: %w", err)
	}

	p.logger.Info("Event batch published",
		zap.Int("batch_size", len(events)))

	return nil
}

// Close closes the Kafka writer
func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

// KafkaConsumer implements event consumption from Kafka
type KafkaConsumer struct {
	reader   *kafka.Reader
	logger   *zap.Logger
	handlers map[string]EventHandler
}

// EventHandler defines the interface for handling events
type EventHandler func(ctx context.Context, event event.Event) error

// NewKafkaConsumer creates a new Kafka event consumer
func NewKafkaConsumer(brokers []string, topic, groupID string, logger *zap.Logger) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
		Logger:      kafka.LoggerFunc(logger.Sugar().Debugf),
		ErrorLogger: kafka.LoggerFunc(logger.Sugar().Errorf),
	})

	return &KafkaConsumer{
		reader:   reader,
		logger:   logger,
		handlers: make(map[string]EventHandler),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (c *KafkaConsumer) RegisterHandler(eventType string, handler EventHandler) {
	c.handlers[eventType] = handler
}

// Start starts consuming events
func (c *KafkaConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.logger.Error("Failed to fetch message", zap.Error(err))
				continue
			}

			// Parse event
			var evt event.Event
			if err := json.Unmarshal(message.Value, &evt); err != nil {
				c.logger.Error("Failed to unmarshal event",
					zap.String("offset", fmt.Sprintf("%d", message.Offset)),
					zap.Error(err))
				// Commit anyway to avoid reprocessing
				if err := c.reader.CommitMessages(ctx, message); err != nil {
					c.logger.Error("Failed to commit message", zap.Error(err))
				}
				continue
			}

			// Find handler
			handler, exists := c.handlers[evt.EventType]
			if !exists {
				c.logger.Warn("No handler registered for event type",
					zap.String("event_type", evt.EventType))
				// Commit anyway
				if err := c.reader.CommitMessages(ctx, message); err != nil {
					c.logger.Error("Failed to commit message", zap.Error(err))
				}
				continue
			}

			// Handle event
			if err := handler(ctx, evt); err != nil {
				c.logger.Error("Failed to handle event",
					zap.String("event_id", evt.ID),
					zap.String("event_type", evt.EventType),
					zap.Error(err))
				// Don't commit on error - will retry
				continue
			}

			// Commit message
			if err := c.reader.CommitMessages(ctx, message); err != nil {
				c.logger.Error("Failed to commit message", zap.Error(err))
			}
		}
	}
}

// Close closes the Kafka reader
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

// ObjectTypeEventPublisher publishes object type related events
type ObjectTypeEventPublisher struct {
	publisher *KafkaPublisher
}

// NewObjectTypeEventPublisher creates a new object type event publisher
func NewObjectTypeEventPublisher(publisher *KafkaPublisher) *ObjectTypeEventPublisher {
	return &ObjectTypeEventPublisher{
		publisher: publisher,
	}
}

// PublishCreated publishes an object type created event
func (p *ObjectTypeEventPublisher) PublishCreated(ctx context.Context, objectTypeID, userID string, data interface{}) error {
	evt := event.Event{
		ID:            generateEventID(),
		EventType:     "object_type.created",
		AggregateID:   objectTypeID,
		AggregateType: "object_type",
		Version:       1,
		Timestamp:     time.Now(),
		UserID:        userID,
		Data:          data,
	}

	return p.publisher.Publish(ctx, evt)
}

// PublishUpdated publishes an object type updated event
func (p *ObjectTypeEventPublisher) PublishUpdated(ctx context.Context, objectTypeID, userID string, version int, data interface{}) error {
	evt := event.Event{
		ID:            generateEventID(),
		EventType:     "object_type.updated",
		AggregateID:   objectTypeID,
		AggregateType: "object_type",
		Version:       version,
		Timestamp:     time.Now(),
		UserID:        userID,
		Data:          data,
	}

	return p.publisher.Publish(ctx, evt)
}

// PublishDeleted publishes an object type deleted event
func (p *ObjectTypeEventPublisher) PublishDeleted(ctx context.Context, objectTypeID, userID string, version int) error {
	evt := event.Event{
		ID:            generateEventID(),
		EventType:     "object_type.deleted",
		AggregateID:   objectTypeID,
		AggregateType: "object_type",
		Version:       version,
		Timestamp:     time.Now(),
		UserID:        userID,
		Data:          map[string]interface{}{"deleted": true},
	}

	return p.publisher.Publish(ctx, evt)
}

// LinkTypeEventPublisher publishes link type related events
type LinkTypeEventPublisher struct {
	publisher *KafkaPublisher
}

// NewLinkTypeEventPublisher creates a new link type event publisher
func NewLinkTypeEventPublisher(publisher *KafkaPublisher) *LinkTypeEventPublisher {
	return &LinkTypeEventPublisher{
		publisher: publisher,
	}
}

// PublishCreated publishes a link type created event
func (p *LinkTypeEventPublisher) PublishCreated(ctx context.Context, linkTypeID, userID string, data interface{}) error {
	evt := event.Event{
		ID:            generateEventID(),
		EventType:     "link_type.created",
		AggregateID:   linkTypeID,
		AggregateType: "link_type",
		Version:       1,
		Timestamp:     time.Now(),
		UserID:        userID,
		Data:          data,
	}

	return p.publisher.Publish(ctx, evt)
}

// PublishUpdated publishes a link type updated event
func (p *LinkTypeEventPublisher) PublishUpdated(ctx context.Context, linkTypeID, userID string, version int, data interface{}) error {
	evt := event.Event{
		ID:            generateEventID(),
		EventType:     "link_type.updated",
		AggregateID:   linkTypeID,
		AggregateType: "link_type",
		Version:       version,
		Timestamp:     time.Now(),
		UserID:        userID,
		Data:          data,
	}

	return p.publisher.Publish(ctx, evt)
}

// PublishDeleted publishes a link type deleted event
func (p *LinkTypeEventPublisher) PublishDeleted(ctx context.Context, linkTypeID, userID string, version int) error {
	evt := event.Event{
		ID:            generateEventID(),
		EventType:     "link_type.deleted",
		AggregateID:   linkTypeID,
		AggregateType: "link_type",
		Version:       version,
		Timestamp:     time.Now(),
		UserID:        userID,
		Data:          map[string]interface{}{"deleted": true},
	}

	return p.publisher.Publish(ctx, evt)
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d_%s", time.Now().UnixNano(), generateRandomString(8))
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}