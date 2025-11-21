package resolvers

import (
	"context"
	"fmt"
	"sync"
)

// ConflictResolver defines the interface for conflict resolution strategies
type ConflictResolver interface {
	// ID returns the unique identifier of the resolver
	ID() string

	// Description returns a human-readable description of the resolver
	Description() string

	// Resolve resolves conflicts between existing and new content
	Resolve(ctx context.Context, input *ConflictInput) ([]byte, error)
}

// ConflictInput contains the input data for conflict resolution
type ConflictInput struct {
	ExistingContent []byte // Existing content (potentially modified by user)
	NewContent      []byte // New content from generator
	FilePath        string // File path for context
}

// ResolverRegistry manages the registration and retrieval of conflict resolvers
type ResolverRegistry struct {
	mu              sync.RWMutex
	resolvers       map[string]ConflictResolver
	defaultResolver ConflictResolver
}

// NewResolverRegistry creates a new resolver registry
func NewResolverRegistry() *ResolverRegistry {
	return &ResolverRegistry{
		resolvers: make(map[string]ConflictResolver),
	}
}

// Register registers a new resolver
func (r *ResolverRegistry) Register(resolver ConflictResolver) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := resolver.ID()
	if _, exists := r.resolvers[id]; exists {
		return fmt.Errorf("resolver %s already registered", id)
	}

	r.resolvers[id] = resolver

	// Set the first registered resolver as default
	if r.defaultResolver == nil {
		r.defaultResolver = resolver
	}

	return nil
}

// MustRegister registers a resolver and panics on error
func (r *ResolverRegistry) MustRegister(resolver ConflictResolver) {
	if err := r.Register(resolver); err != nil {
		panic(err)
	}
}

// SetDefault sets the default resolver by ID
func (r *ResolverRegistry) SetDefault(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	resolver, exists := r.resolvers[id]
	if !exists {
		return fmt.Errorf("resolver %s not found", id)
	}

	r.defaultResolver = resolver
	return nil
}

// MustSetDefault sets the default resolver and panics on error
func (r *ResolverRegistry) MustSetDefault(id string) {
	if err := r.SetDefault(id); err != nil {
		panic(err)
	}
}

// GetDefault returns the default resolver
func (r *ResolverRegistry) GetDefault() ConflictResolver {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.defaultResolver == nil {
		panic("no default resolver set")
	}

	return r.defaultResolver
}

// Get retrieves a resolver by ID
func (r *ResolverRegistry) Get(id string) (ConflictResolver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resolver, exists := r.resolvers[id]
	return resolver, exists
}

// MustGet retrieves a resolver by ID and panics if not found
func (r *ResolverRegistry) MustGet(id string) ConflictResolver {
	resolver, exists := r.Get(id)
	if !exists {
		panic(fmt.Sprintf("resolver %s not found", id))
	}
	return resolver
}

// List returns all registered resolvers
func (r *ResolverRegistry) List() []ConflictResolver {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ConflictResolver, 0, len(r.resolvers))
	for _, resolver := range r.resolvers {
		result = append(result, resolver)
	}

	return result
}

// DefaultResolverRegistry is the global default resolver registry
var DefaultResolverRegistry = NewResolverRegistry()
