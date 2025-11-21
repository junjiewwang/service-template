package mergers

import (
	"context"
	"fmt"
	"sync"
)

// ContentMerger defines the interface for content merging strategies
type ContentMerger interface {
	// ID returns the unique identifier of the merger
	ID() string

	// Description returns a human-readable description of the merger
	Description() string

	// Merge merges existing content with new content
	Merge(ctx context.Context, input *MergeInput) ([]byte, error)
}

// MergeInput contains the input data for content merging
type MergeInput struct {
	ExistingContent []byte // Existing file content
	NewContent      []byte // New content to merge
	FilePath        string // File path for context
}

// MergerRegistry manages the registration and retrieval of content mergers
type MergerRegistry struct {
	mu            sync.RWMutex
	mergers       map[string]ContentMerger
	defaultMerger ContentMerger
}

// NewMergerRegistry creates a new merger registry
func NewMergerRegistry() *MergerRegistry {
	return &MergerRegistry{
		mergers: make(map[string]ContentMerger),
	}
}

// Register registers a new merger
func (r *MergerRegistry) Register(merger ContentMerger) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := merger.ID()
	if _, exists := r.mergers[id]; exists {
		return fmt.Errorf("merger %s already registered", id)
	}

	r.mergers[id] = merger

	// Set the first registered merger as default
	if r.defaultMerger == nil {
		r.defaultMerger = merger
	}

	return nil
}

// MustRegister registers a merger and panics on error
func (r *MergerRegistry) MustRegister(merger ContentMerger) {
	if err := r.Register(merger); err != nil {
		panic(err)
	}
}

// SetDefault sets the default merger by ID
func (r *MergerRegistry) SetDefault(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	merger, exists := r.mergers[id]
	if !exists {
		return fmt.Errorf("merger %s not found", id)
	}

	r.defaultMerger = merger
	return nil
}

// MustSetDefault sets the default merger and panics on error
func (r *MergerRegistry) MustSetDefault(id string) {
	if err := r.SetDefault(id); err != nil {
		panic(err)
	}
}

// GetDefault returns the default merger
func (r *MergerRegistry) GetDefault() ContentMerger {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.defaultMerger == nil {
		panic("no default merger set")
	}

	return r.defaultMerger
}

// Get retrieves a merger by ID
func (r *MergerRegistry) Get(id string) (ContentMerger, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	merger, exists := r.mergers[id]
	return merger, exists
}

// MustGet retrieves a merger by ID and panics if not found
func (r *MergerRegistry) MustGet(id string) ContentMerger {
	merger, exists := r.Get(id)
	if !exists {
		panic(fmt.Sprintf("merger %s not found", id))
	}
	return merger
}

// List returns all registered mergers
func (r *MergerRegistry) List() []ContentMerger {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ContentMerger, 0, len(r.mergers))
	for _, merger := range r.mergers {
		result = append(result, merger)
	}

	return result
}

// DefaultMergerRegistry is the global default merger registry
var DefaultMergerRegistry = NewMergerRegistry()
