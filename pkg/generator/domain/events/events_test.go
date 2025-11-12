package events

import (
	"testing"
	"time"
)

func TestBaseEvent(t *testing.T) {
	data := map[string]interface{}{
		"key": "value",
	}

	event := NewBaseEvent("test.event", "test-aggregate", data)

	if event.EventName() != "test.event" {
		t.Errorf("EventName() = %v, want test.event", event.EventName())
	}

	if event.AggregateID() != "test-aggregate" {
		t.Errorf("AggregateID() = %v, want test-aggregate", event.AggregateID())
	}

	if event.EventData()["key"] != "value" {
		t.Error("EventData() did not return correct data")
	}

	if time.Since(event.OccurredAt()) > time.Second {
		t.Error("OccurredAt() is not recent")
	}
}

func TestSimpleEventPublisher(t *testing.T) {
	publisher := NewSimpleEventPublisher()

	if publisher.HandlerCount() != 0 {
		t.Errorf("HandlerCount() = %v, want 0", publisher.HandlerCount())
	}

	handler := NewMetricsEventHandler()
	publisher.Subscribe(handler)

	if publisher.HandlerCount() != 1 {
		t.Errorf("HandlerCount() = %v, want 1", publisher.HandlerCount())
	}

	event := NewGenerationStartedEvent("test-generator", map[string]interface{}{})

	err := publisher.Publish(event)
	if err != nil {
		t.Errorf("Publish() error = %v", err)
	}

	if handler.GetCount(EventGenerationStarted) != 1 {
		t.Errorf("Event count = %v, want 1", handler.GetCount(EventGenerationStarted))
	}

	publisher.Unsubscribe(handler)

	if publisher.HandlerCount() != 0 {
		t.Errorf("HandlerCount() = %v, want 0 after unsubscribe", publisher.HandlerCount())
	}
}

func TestLoggingEventHandler(t *testing.T) {
	handler := NewLoggingEventHandler("TEST")

	if !handler.CanHandle("any.event") {
		t.Error("CanHandle() should return true for all events")
	}

	event := NewGenerationStartedEvent("test-generator", map[string]interface{}{
		"test": "data",
	})

	err := handler.Handle(event)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
}

func TestMetricsEventHandler(t *testing.T) {
	handler := NewMetricsEventHandler()

	event1 := NewGenerationStartedEvent("gen1", nil)
	event2 := NewGenerationCompletedEvent("gen1", nil)
	event3 := NewGenerationStartedEvent("gen2", nil)

	handler.Handle(event1)
	handler.Handle(event2)
	handler.Handle(event3)

	if handler.GetCount(EventGenerationStarted) != 2 {
		t.Errorf("GetCount(started) = %v, want 2", handler.GetCount(EventGenerationStarted))
	}

	if handler.GetCount(EventGenerationCompleted) != 1 {
		t.Errorf("GetCount(completed) = %v, want 1", handler.GetCount(EventGenerationCompleted))
	}

	metrics := handler.GetMetrics()
	if len(metrics) != 2 {
		t.Errorf("GetMetrics() returned %v events, want 2", len(metrics))
	}

	handler.Reset()
	if handler.GetCount(EventGenerationStarted) != 0 {
		t.Error("Reset() did not clear metrics")
	}
}

func TestFilteredEventHandler(t *testing.T) {
	metricsHandler := NewMetricsEventHandler()
	filteredHandler := NewFilteredEventHandler(metricsHandler, "generation")

	if !filteredHandler.CanHandle("generation.started") {
		t.Error("CanHandle() should return true for generation events")
	}

	if filteredHandler.CanHandle("validation.started") {
		t.Error("CanHandle() should return false for validation events")
	}

	event1 := NewGenerationStartedEvent("gen1", nil)
	event2 := NewValidationStartedEvent("gen1", nil)

	filteredHandler.Handle(event1)
	filteredHandler.Handle(event2)

	// Only generation event should be counted
	if metricsHandler.GetCount(EventGenerationStarted) != 1 {
		t.Errorf("Filtered handler counted wrong events")
	}
}

func TestCompositeEventHandler(t *testing.T) {
	handler1 := NewMetricsEventHandler()
	handler2 := NewMetricsEventHandler()

	composite := NewCompositeEventHandler(handler1, handler2)

	if !composite.CanHandle("any.event") {
		t.Error("CanHandle() should return true")
	}

	event := NewGenerationStartedEvent("gen1", nil)

	err := composite.Handle(event)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}

	// Both handlers should have counted the event
	if handler1.GetCount(EventGenerationStarted) != 1 {
		t.Error("Handler1 did not count event")
	}

	if handler2.GetCount(EventGenerationStarted) != 1 {
		t.Error("Handler2 did not count event")
	}

	// Test AddHandler
	handler3 := NewMetricsEventHandler()
	composite.AddHandler(handler3)

	composite.Handle(event)

	if handler3.GetCount(EventGenerationStarted) != 1 {
		t.Error("Newly added handler did not count event")
	}
}

func TestGeneratorEvents(t *testing.T) {
	tests := []struct {
		name      string
		eventFunc func() Event
		eventName string
	}{
		{
			name: "GenerationStarted",
			eventFunc: func() Event {
				return NewGenerationStartedEvent("test", nil)
			},
			eventName: EventGenerationStarted,
		},
		{
			name: "GenerationCompleted",
			eventFunc: func() Event {
				return NewGenerationCompletedEvent("test", nil)
			},
			eventName: EventGenerationCompleted,
		},
		{
			name: "GenerationFailed",
			eventFunc: func() Event {
				return NewGenerationFailedEvent("test", nil)
			},
			eventName: EventGenerationFailed,
		},
		{
			name: "ValidationStarted",
			eventFunc: func() Event {
				return NewValidationStartedEvent("test", nil)
			},
			eventName: EventValidationStarted,
		},
		{
			name: "ValidationCompleted",
			eventFunc: func() Event {
				return NewValidationCompletedEvent("test", nil)
			},
			eventName: EventValidationCompleted,
		},
		{
			name: "ValidationFailed",
			eventFunc: func() Event {
				return NewValidationFailedEvent("test", nil)
			},
			eventName: EventValidationFailed,
		},
		{
			name: "TemplateRendered",
			eventFunc: func() Event {
				return NewTemplateRenderedEvent("test", nil)
			},
			eventName: EventTemplateRendered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.eventFunc()

			if event.EventName() != tt.eventName {
				t.Errorf("EventName() = %v, want %v", event.EventName(), tt.eventName)
			}

			if event.AggregateID() != "test" {
				t.Errorf("AggregateID() = %v, want test", event.AggregateID())
			}
		})
	}
}

func TestPublisherConcurrency(t *testing.T) {
	publisher := NewSimpleEventPublisher()
	handler := NewMetricsEventHandler()
	publisher.Subscribe(handler)

	// Publish events concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			event := NewGenerationStartedEvent("test", nil)
			publisher.Publish(event)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	if handler.GetCount(EventGenerationStarted) != 10 {
		t.Errorf("Concurrent publish failed, got %v events, want 10",
			handler.GetCount(EventGenerationStarted))
	}
}
