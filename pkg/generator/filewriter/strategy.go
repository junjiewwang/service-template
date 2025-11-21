package filewriter

import (
	"context"
	"fmt"
	"sync"
)

// WriteStrategy defines the interface for file writing strategies
type WriteStrategy interface {
	// ID returns the unique identifier of the strategy
	ID() string

	// Description returns a human-readable description of the strategy
	Description() string

	// Write executes the file writing operation
	Write(ctx context.Context, path string, content []byte) error
}

// StrategyRegistry manages the registration and retrieval of write strategies
type StrategyRegistry struct {
	mu              sync.RWMutex
	strategies      map[string]WriteStrategy
	defaultStrategy WriteStrategy
}

// NewStrategyRegistry creates a new strategy registry
func NewStrategyRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategies: make(map[string]WriteStrategy),
	}
}

// Register registers a new strategy
func (r *StrategyRegistry) Register(strategy WriteStrategy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := strategy.ID()
	if _, exists := r.strategies[id]; exists {
		return fmt.Errorf("strategy %s already registered", id)
	}

	r.strategies[id] = strategy

	// Set the first registered strategy as default
	if r.defaultStrategy == nil {
		r.defaultStrategy = strategy
	}

	return nil
}

// MustRegister registers a strategy and panics on error
func (r *StrategyRegistry) MustRegister(strategy WriteStrategy) {
	if err := r.Register(strategy); err != nil {
		panic(err)
	}
}

// SetDefault sets the default strategy by ID
func (r *StrategyRegistry) SetDefault(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	strategy, exists := r.strategies[id]
	if !exists {
		return fmt.Errorf("strategy %s not found", id)
	}

	r.defaultStrategy = strategy
	return nil
}

// MustSetDefault sets the default strategy and panics on error
func (r *StrategyRegistry) MustSetDefault(id string) {
	if err := r.SetDefault(id); err != nil {
		panic(err)
	}
}

// GetDefault returns the default strategy
func (r *StrategyRegistry) GetDefault() WriteStrategy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.defaultStrategy == nil {
		panic("no default strategy set")
	}

	return r.defaultStrategy
}

// Get retrieves a strategy by ID
func (r *StrategyRegistry) Get(id string) (WriteStrategy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	strategy, exists := r.strategies[id]
	return strategy, exists
}

// MustGet retrieves a strategy by ID and panics if not found
func (r *StrategyRegistry) MustGet(id string) WriteStrategy {
	strategy, exists := r.Get(id)
	if !exists {
		panic(fmt.Sprintf("strategy %s not found", id))
	}
	return strategy
}

// List returns all registered strategies
func (r *StrategyRegistry) List() []WriteStrategy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]WriteStrategy, 0, len(r.strategies))
	for _, strategy := range r.strategies {
		result = append(result, strategy)
	}

	return result
}

// DefaultStrategyRegistry is the global default strategy registry
var DefaultStrategyRegistry = NewStrategyRegistry()
