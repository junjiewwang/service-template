package repositories

import (
	"sync"
)

// EmbeddedTemplateRepository stores templates in memory (embedded in binary)
type EmbeddedTemplateRepository struct {
	templates map[string]*Template
	mu        sync.RWMutex
}

// NewEmbeddedTemplateRepository creates a new embedded template repository
func NewEmbeddedTemplateRepository() *EmbeddedTemplateRepository {
	return &EmbeddedTemplateRepository{
		templates: make(map[string]*Template),
	}
}

// Register registers a template (used during initialization)
func (r *EmbeddedTemplateRepository) Register(template *Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[template.Name]; exists {
		return &TemplateAlreadyExistsError{Name: template.Name}
	}

	r.templates[template.Name] = template
	return nil
}

// Get retrieves a template by name
func (r *EmbeddedTemplateRepository) Get(name string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.templates[name]
	if !exists {
		return nil, &TemplateNotFoundError{Name: name}
	}

	return template, nil
}

// List returns all available templates
func (r *EmbeddedTemplateRepository) List() ([]*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]*Template, 0, len(r.templates))
	for _, template := range r.templates {
		templates = append(templates, template)
	}

	return templates, nil
}

// ListByCategory returns templates in a specific category
func (r *EmbeddedTemplateRepository) ListByCategory(category string) ([]*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]*Template, 0)
	for _, template := range r.templates {
		if template.Category == category {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

// Exists checks if a template exists
func (r *EmbeddedTemplateRepository) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.templates[name]
	return exists
}

// Save is not supported for embedded repository
func (r *EmbeddedTemplateRepository) Save(template *Template) error {
	return &ReadOnlyRepositoryError{Operation: "Save"}
}

// Delete is not supported for embedded repository
func (r *EmbeddedTemplateRepository) Delete(name string) error {
	return &ReadOnlyRepositoryError{Operation: "Delete"}
}

// Count returns the number of templates
func (r *EmbeddedTemplateRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.templates)
}

// Clear removes all templates (for testing)
func (r *EmbeddedTemplateRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates = make(map[string]*Template)
}
