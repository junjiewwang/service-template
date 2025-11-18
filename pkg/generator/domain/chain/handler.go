package chain

// Handler is the base interface for all handlers in the chain
type Handler interface {
	// Handle processes the request and passes it to the next handler
	Handle(ctx *ProcessingContext) error

	// SetNext sets the next handler in the chain
	SetNext(handler Handler) Handler

	// GetName returns the handler name for logging and debugging
	GetName() string
}

// ParserHandler handles configuration parsing
type ParserHandler interface {
	Handler
	// Parse parses configuration from raw data
	Parse(ctx *ProcessingContext) error
}

// ValidatorHandler handles configuration validation
type ValidatorHandler interface {
	Handler
	// Validate validates the parsed configuration
	Validate(ctx *ProcessingContext) error
}

// GeneratorHandler handles file generation
type GeneratorHandler interface {
	Handler
	// Generate generates files based on configuration
	Generate(ctx *ProcessingContext) error
}

// BaseHandler provides default implementation for Handler interface
type BaseHandler struct {
	next Handler
	name string
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(name string) *BaseHandler {
	return &BaseHandler{
		name: name,
	}
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

// GetName returns the handler name
func (h *BaseHandler) GetName() string {
	return h.name
}

// CallNext calls the next handler in the chain
func (h *BaseHandler) CallNext(ctx *ProcessingContext) error {
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return nil
}

// Handle is the default implementation (should be overridden)
func (h *BaseHandler) Handle(ctx *ProcessingContext) error {
	return h.CallNext(ctx)
}
