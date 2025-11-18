package chain

import (
	"context"
	"sync"
)

// ProcessingContext holds the context for processing through the chain
type ProcessingContext struct {
	// Context for cancellation and timeout
	Context context.Context

	// Raw configuration data
	RawConfig map[string]interface{}

	// Parsed domain models (keyed by domain name)
	DomainModels sync.Map

	// Validation errors (keyed by domain name)
	ValidationErrors sync.Map

	// Generated files (keyed by file path)
	GeneratedFiles sync.Map

	// Metadata for tracking
	Metadata map[string]interface{}

	// Error aggregation
	errors []error
	mu     sync.RWMutex
}

// NewProcessingContext creates a new processing context
func NewProcessingContext(ctx context.Context, rawConfig map[string]interface{}) *ProcessingContext {
	return &ProcessingContext{
		Context:   ctx,
		RawConfig: rawConfig,
		Metadata:  make(map[string]interface{}),
		errors:    make([]error, 0),
	}
}

// SetDomainModel stores a parsed domain model
func (c *ProcessingContext) SetDomainModel(domain string, model interface{}) {
	c.DomainModels.Store(domain, model)
}

// GetDomainModel retrieves a parsed domain model
func (c *ProcessingContext) GetDomainModel(domain string) (interface{}, bool) {
	return c.DomainModels.Load(domain)
}

// AddValidationError adds a validation error for a domain
func (c *ProcessingContext) AddValidationError(domain string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	errors, _ := c.ValidationErrors.LoadOrStore(domain, []error{})
	errorList := errors.([]error)
	errorList = append(errorList, err)
	c.ValidationErrors.Store(domain, errorList)
	c.errors = append(c.errors, err)
}

// GetValidationErrors retrieves validation errors for a domain
func (c *ProcessingContext) GetValidationErrors(domain string) []error {
	errors, ok := c.ValidationErrors.Load(domain)
	if !ok {
		return nil
	}
	return errors.([]error)
}

// AddGeneratedFile records a generated file
func (c *ProcessingContext) AddGeneratedFile(path string, content []byte) {
	c.GeneratedFiles.Store(path, content)
}

// GetGeneratedFile retrieves a generated file
func (c *ProcessingContext) GetGeneratedFile(path string) ([]byte, bool) {
	content, ok := c.GeneratedFiles.Load(path)
	if !ok {
		return nil, false
	}
	return content.([]byte), true
}

// AddError adds a general error
func (c *ProcessingContext) AddError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors = append(c.errors, err)
}

// GetErrors returns all errors
func (c *ProcessingContext) GetErrors() []error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]error{}, c.errors...)
}

// HasErrors checks if there are any errors
func (c *ProcessingContext) HasErrors() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.errors) > 0
}

// SetMetadata sets a metadata value
func (c *ProcessingContext) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Metadata[key] = value
}

// GetMetadata retrieves a metadata value
func (c *ProcessingContext) GetMetadata(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.Metadata[key]
	return value, ok
}

// HasValidationErrors checks if there are any validation errors
func (c *ProcessingContext) HasValidationErrors() bool {
	hasErrors := false
	c.ValidationErrors.Range(func(key, value interface{}) bool {
		if errors, ok := value.([]error); ok && len(errors) > 0 {
			hasErrors = true
			return false // stop iteration
		}
		return true
	})
	return hasErrors
}

// GetAllValidationErrors returns all validation errors grouped by domain
func (c *ProcessingContext) GetAllValidationErrors() map[string][]error {
	result := make(map[string][]error)
	c.ValidationErrors.Range(func(key, value interface{}) bool {
		if domain, ok := key.(string); ok {
			if errors, ok := value.([]error); ok {
				result[domain] = errors
			}
		}
		return true
	})
	return result
}

// GetAllGeneratedFiles returns all generated files
func (c *ProcessingContext) GetAllGeneratedFiles() map[string][]byte {
	result := make(map[string][]byte)
	c.GeneratedFiles.Range(func(key, value interface{}) bool {
		if path, ok := key.(string); ok {
			if content, ok := value.([]byte); ok {
				result[path] = content
			}
		}
		return true
	})
	return result
}
