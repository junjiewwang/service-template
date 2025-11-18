package chain

import (
	"fmt"
	"sort"
	"sync"
)

// DomainRegistry manages domain factory registration and chain building
type DomainRegistry struct {
	factories map[string]DomainFactory
	mu        sync.RWMutex
}

// NewDomainRegistry creates a new domain registry
func NewDomainRegistry() *DomainRegistry {
	return &DomainRegistry{
		factories: make(map[string]DomainFactory),
	}
}

// Register registers a domain factory
func (r *DomainRegistry) Register(factory DomainFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := factory.GetName()
	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("domain factory %s already registered", name)
	}

	r.factories[name] = factory
	return nil
}

// RegisterAll registers multiple domain factories
func (r *DomainRegistry) RegisterAll(factories ...DomainFactory) error {
	for _, factory := range factories {
		if err := r.Register(factory); err != nil {
			return err
		}
	}
	return nil
}

// Unregister unregisters a domain factory
func (r *DomainRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.factories, name)
}

// Get retrieves a domain factory by name
func (r *DomainRegistry) Get(name string) (DomainFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, ok := r.factories[name]
	return factory, ok
}

// GetAll returns all registered factories
func (r *DomainRegistry) GetAll() []DomainFactory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factories := make([]DomainFactory, 0, len(r.factories))
	for _, factory := range r.factories {
		factories = append(factories, factory)
	}
	return factories
}

// GetEnabled returns all enabled factories
func (r *DomainRegistry) GetEnabled() []DomainFactory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factories := make([]DomainFactory, 0)
	for _, factory := range r.factories {
		if factory.IsEnabled() {
			factories = append(factories, factory)
		}
	}
	return factories
}

// BuildParseChain builds a parser chain from registered factories
func (r *DomainRegistry) BuildParseChain() Handler {
	factories := r.getSortedFactories()

	builder := NewChainBuilder()
	for _, factory := range factories {
		if handler := factory.CreateParserHandler(); handler != nil {
			builder.Add(handler)
		}
	}

	return builder.BuildWithLogging()
}

// BuildValidateChain builds a validator chain from registered factories
func (r *DomainRegistry) BuildValidateChain() Handler {
	factories := r.getSortedFactories()

	builder := NewChainBuilder()
	for _, factory := range factories {
		if handler := factory.CreateValidatorHandler(); handler != nil {
			builder.Add(handler)
		}
	}

	return builder.BuildWithLogging()
}

// BuildGenerateChain builds a generator chain from registered factories
func (r *DomainRegistry) BuildGenerateChain() Handler {
	factories := r.getSortedFactories()

	builder := NewChainBuilder()
	for _, factory := range factories {
		if handler := factory.CreateGeneratorHandler(); handler != nil {
			builder.Add(handler)
		}
	}

	return builder.BuildWithLogging()
}

// getSortedFactories returns enabled factories sorted by priority
func (r *DomainRegistry) getSortedFactories() []DomainFactory {
	factories := r.GetEnabled()

	// Sort by priority (lower value = higher priority)
	sort.Slice(factories, func(i, j int) bool {
		return factories[i].GetPriority() < factories[j].GetPriority()
	})

	return factories
}

// Clear removes all registered factories
func (r *DomainRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories = make(map[string]DomainFactory)
}

// Count returns the number of registered factories
func (r *DomainRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.factories)
}
