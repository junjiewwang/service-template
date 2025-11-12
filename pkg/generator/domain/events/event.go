package events

import (
	"time"
)

// Event represents a domain event
type Event interface {
	// EventName returns the name of the event
	EventName() string

	// OccurredAt returns when the event occurred
	OccurredAt() time.Time

	// AggregateID returns the ID of the aggregate that produced the event
	AggregateID() string

	// EventData returns the event data
	EventData() map[string]interface{}
}

// BaseEvent provides common event functionality
type BaseEvent struct {
	name        string
	occurredAt  time.Time
	aggregateID string
	data        map[string]interface{}
}

// NewBaseEvent creates a new base event
func NewBaseEvent(name, aggregateID string, data map[string]interface{}) BaseEvent {
	return BaseEvent{
		name:        name,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
		data:        data,
	}
}

// EventName returns the event name
func (e BaseEvent) EventName() string {
	return e.name
}

// OccurredAt returns when the event occurred
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID returns the aggregate ID
func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

// EventData returns the event data
func (e BaseEvent) EventData() map[string]interface{} {
	return e.data
}

// EventHandler handles domain events
type EventHandler interface {
	// Handle processes an event
	Handle(event Event) error

	// CanHandle returns true if this handler can handle the event
	CanHandle(eventName string) bool
}

// EventPublisher publishes domain events
type EventPublisher interface {
	// Publish publishes an event to all registered handlers
	Publish(event Event) error

	// Subscribe registers an event handler
	Subscribe(handler EventHandler)

	// Unsubscribe removes an event handler
	Unsubscribe(handler EventHandler)
}
