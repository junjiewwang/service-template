package events

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

// LoggingEventHandler logs all events
type LoggingEventHandler struct {
	prefix string
}

// NewLoggingEventHandler creates a new logging event handler
func NewLoggingEventHandler(prefix string) *LoggingEventHandler {
	return &LoggingEventHandler{
		prefix: prefix,
	}
}

// Handle logs the event
func (h *LoggingEventHandler) Handle(event Event) error {
	log.Printf("[%s] Event: %s | Aggregate: %s | Time: %s | Data: %v",
		h.prefix,
		event.EventName(),
		event.AggregateID(),
		event.OccurredAt().Format("2006-01-02 15:04:05"),
		event.EventData(),
	)
	return nil
}

// CanHandle returns true for all events
func (h *LoggingEventHandler) CanHandle(eventName string) bool {
	return true
}

// MetricsEventHandler collects metrics from events
type MetricsEventHandler struct {
	metrics map[string]int
	mu      sync.RWMutex
}

// NewMetricsEventHandler creates a new metrics event handler
func NewMetricsEventHandler() *MetricsEventHandler {
	return &MetricsEventHandler{
		metrics: make(map[string]int),
	}
}

// Handle increments the counter for the event
func (h *MetricsEventHandler) Handle(event Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.metrics[event.EventName()]++
	return nil
}

// CanHandle returns true for all events
func (h *MetricsEventHandler) CanHandle(eventName string) bool {
	return true
}

// GetMetrics returns the collected metrics
func (h *MetricsEventHandler) GetMetrics() map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make(map[string]int)
	for k, v := range h.metrics {
		result[k] = v
	}
	return result
}

// GetCount returns the count for a specific event
func (h *MetricsEventHandler) GetCount(eventName string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.metrics[eventName]
}

// Reset resets all metrics
func (h *MetricsEventHandler) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.metrics = make(map[string]int)
}

// FilteredEventHandler filters events by name pattern
type FilteredEventHandler struct {
	patterns []string
	handler  EventHandler
}

// NewFilteredEventHandler creates a new filtered event handler
func NewFilteredEventHandler(handler EventHandler, patterns ...string) *FilteredEventHandler {
	return &FilteredEventHandler{
		patterns: patterns,
		handler:  handler,
	}
}

// Handle delegates to the wrapped handler if the event matches
func (h *FilteredEventHandler) Handle(event Event) error {
	return h.handler.Handle(event)
}

// CanHandle returns true if the event name matches any pattern
func (h *FilteredEventHandler) CanHandle(eventName string) bool {
	for _, pattern := range h.patterns {
		if strings.Contains(eventName, pattern) {
			return h.handler.CanHandle(eventName)
		}
	}
	return false
}

// CompositeEventHandler combines multiple handlers
type CompositeEventHandler struct {
	handlers []EventHandler
}

// NewCompositeEventHandler creates a new composite event handler
func NewCompositeEventHandler(handlers ...EventHandler) *CompositeEventHandler {
	return &CompositeEventHandler{
		handlers: handlers,
	}
}

// Handle calls all handlers
func (h *CompositeEventHandler) Handle(event Event) error {
	var errors []error
	for _, handler := range h.handlers {
		if handler.CanHandle(event.EventName()) {
			if err := handler.Handle(event); err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("composite handler errors: %v", errors)
	}

	return nil
}

// CanHandle returns true if any handler can handle the event
func (h *CompositeEventHandler) CanHandle(eventName string) bool {
	for _, handler := range h.handlers {
		if handler.CanHandle(eventName) {
			return true
		}
	}
	return false
}

// AddHandler adds a handler to the composite
func (h *CompositeEventHandler) AddHandler(handler EventHandler) {
	h.handlers = append(h.handlers, handler)
}
