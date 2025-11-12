package core

import (
	"fmt"
	"sync"
)

// Registry manages generator registration and creation
type Registry struct {
	creators map[string]GeneratorCreator
	mu       sync.RWMutex
}

// NewRegistry creates a new generator registry
func NewRegistry() *Registry {
	return &Registry{
		creators: make(map[string]GeneratorCreator),
	}
}

// Register registers a generator creator with a unique type identifier
func (r *Registry) Register(generatorType string, creator GeneratorCreator) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.creators[generatorType]; exists {
		panic(fmt.Sprintf("generator type %s is already registered", generatorType))
	}
	r.creators[generatorType] = creator
}

// Get retrieves a generator creator by type
func (r *Registry) Get(generatorType string) (GeneratorCreator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	creator, exists := r.creators[generatorType]
	return creator, exists
}

// GetAll returns all registered generator types
func (r *Registry) GetAll() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.creators))
	for t := range r.creators {
		types = append(types, t)
	}
	return types
}

// DefaultRegistry is the global registry instance
var DefaultRegistry = NewRegistry()
