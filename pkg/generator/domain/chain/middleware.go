package chain

import (
	"fmt"
	"log"
	"time"
)

// Middleware wraps a handler with additional functionality
type Middleware func(Handler) Handler

// LoggingMiddleware wraps a handler with logging
type LoggingMiddleware struct {
	*BaseHandler
	wrapped Handler
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(handler Handler) Handler {
	return &LoggingMiddleware{
		BaseHandler: NewBaseHandler(fmt.Sprintf("logging(%s)", handler.GetName())),
		wrapped:     handler,
	}
}

// Handle executes the wrapped handler with logging
func (m *LoggingMiddleware) Handle(ctx *ProcessingContext) error {
	start := time.Now()
	log.Printf("[%s] Starting handler: %s", m.wrapped.GetName(), m.wrapped.GetName())

	// Execute the wrapped handler (which will call its own CallNext internally)
	err := m.wrapped.Handle(ctx)

	duration := time.Since(start)
	if err != nil {
		log.Printf("[%s] Handler failed: %s (duration: %v, error: %v)",
			m.wrapped.GetName(), m.wrapped.GetName(), duration, err)
	} else {
		log.Printf("[%s] Handler completed: %s (duration: %v)",
			m.wrapped.GetName(), m.wrapped.GetName(), duration)
	}

	// Note: We don't call m.CallNext() here because the wrapped handler
	// already handles the chain continuation internally
	return err
}

// MetricsMiddleware wraps a handler with metrics collection
type MetricsMiddleware struct {
	*BaseHandler
	wrapped Handler
	metrics *HandlerMetrics
}

// HandlerMetrics holds metrics for a handler
type HandlerMetrics struct {
	HandlerName    string
	ExecutionCount int64
	TotalDuration  time.Duration
	ErrorCount     int64
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(handler Handler) Handler {
	return &MetricsMiddleware{
		BaseHandler: NewBaseHandler(fmt.Sprintf("metrics(%s)", handler.GetName())),
		wrapped:     handler,
		metrics: &HandlerMetrics{
			HandlerName: handler.GetName(),
		},
	}
}

// Handle executes the wrapped handler with metrics collection
func (m *MetricsMiddleware) Handle(ctx *ProcessingContext) error {
	start := time.Now()
	m.metrics.ExecutionCount++

	err := m.wrapped.Handle(ctx)

	duration := time.Since(start)
	m.metrics.TotalDuration += duration

	if err != nil {
		m.metrics.ErrorCount++
	}

	return err
}

// GetMetrics returns the collected metrics
func (m *MetricsMiddleware) GetMetrics() *HandlerMetrics {
	return m.metrics
}

// RecoveryMiddleware wraps a handler with panic recovery
type RecoveryMiddleware struct {
	*BaseHandler
	wrapped Handler
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(handler Handler) Handler {
	return &RecoveryMiddleware{
		BaseHandler: NewBaseHandler(fmt.Sprintf("recovery(%s)", handler.GetName())),
		wrapped:     handler,
	}
}

// Handle executes the wrapped handler with panic recovery
func (m *RecoveryMiddleware) Handle(ctx *ProcessingContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in handler %s: %v", m.wrapped.GetName(), r)
			log.Printf("[PANIC] Handler %s panicked: %v", m.wrapped.GetName(), r)
		}
	}()

	return m.wrapped.Handle(ctx)
}
