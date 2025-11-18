package chain

// ChainBuilder builds a chain of handlers
type ChainBuilder struct {
	handlers []Handler
}

// NewChainBuilder creates a new chain builder
func NewChainBuilder() *ChainBuilder {
	return &ChainBuilder{
		handlers: make([]Handler, 0),
	}
}

// Add adds a handler to the chain
func (b *ChainBuilder) Add(handler Handler) *ChainBuilder {
	b.handlers = append(b.handlers, handler)
	return b
}

// Build builds the chain and returns the first handler
func (b *ChainBuilder) Build() Handler {
	if len(b.handlers) == 0 {
		return nil
	}

	// Link handlers together
	for i := 0; i < len(b.handlers)-1; i++ {
		b.handlers[i].SetNext(b.handlers[i+1])
	}

	return b.handlers[0]
}

// BuildWithLogging builds the chain with logging middleware
func (b *ChainBuilder) BuildWithLogging() Handler {
	if len(b.handlers) == 0 {
		return nil
	}

	// First, link the original handlers together
	for i := 0; i < len(b.handlers)-1; i++ {
		b.handlers[i].SetNext(b.handlers[i+1])
	}

	// Then wrap each handler with logging middleware
	// Note: We wrap them individually so each handler's logging is independent
	wrappedHandlers := make([]Handler, len(b.handlers))
	for i, handler := range b.handlers {
		wrappedHandlers[i] = NewLoggingMiddleware(handler)
	}

	// Link the wrapped handlers
	for i := 0; i < len(wrappedHandlers)-1; i++ {
		wrappedHandlers[i].SetNext(wrappedHandlers[i+1])
	}

	return wrappedHandlers[0]
}

// BuildWithMiddleware builds the chain with custom middleware
func (b *ChainBuilder) BuildWithMiddleware(middlewareFunc func(Handler) Handler) Handler {
	if len(b.handlers) == 0 {
		return nil
	}

	// Wrap each handler with middleware
	wrappedHandlers := make([]Handler, len(b.handlers))
	for i, handler := range b.handlers {
		wrappedHandlers[i] = middlewareFunc(handler)
	}

	// Link wrapped handlers together
	for i := 0; i < len(wrappedHandlers)-1; i++ {
		wrappedHandlers[i].SetNext(wrappedHandlers[i+1])
	}

	return wrappedHandlers[0]
}
