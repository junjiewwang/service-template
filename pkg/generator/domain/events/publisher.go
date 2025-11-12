package events

import (
	"fmt"
	"sync"
)

// SimpleEventPublisher is a simple in-memory event publisher
type SimpleEventPublisher struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

// NewSimpleEventPublisher creates a new simple event publisher
func NewSimpleEventPublisher() *SimpleEventPublisher {
	return &SimpleEventPublisher{
		handlers: make([]EventHandler, 0),
	}
}

// Publish publishes an event to all registered handlers
func (p *SimpleEventPublisher) Publish(event Event) error {
	p.mu.RLock()
	handlers := make([]EventHandler, len(p.handlers))
	copy(handlers, p.handlers)
	p.mu.RUnlock()

	var errors []error
	for _, handler := range handlers {
		if handler.CanHandle(event.EventName()) {
			if err := handler.Handle(event); err != nil {
				errors = append(errors, fmt.Errorf("handler error: %w", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("event publishing errors: %v", errors)
	}

	return nil
}

// Subscribe registers an event handler
func (p *SimpleEventPublisher) Subscribe(handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = append(p.handlers, handler)
}

// Unsubscribe removes an event handler
func (p *SimpleEventPublisher) Unsubscribe(handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, h := range p.handlers {
		if h == handler {
			p.handlers = append(p.handlers[:i], p.handlers[i+1:]...)
			break
		}
	}
}

// HandlerCount returns the number of registered handlers
func (p *SimpleEventPublisher) HandlerCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.handlers)
}

// Clear removes all handlers
func (p *SimpleEventPublisher) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = make([]EventHandler, 0)
}
