package repositories

import (
	"fmt"
)

// Template represents a template with its content and metadata
type Template struct {
	Name        string
	Content     string
	Description string
	Category    string
	Tags        []string
}

// TemplateRepository defines the interface for template storage
type TemplateRepository interface {
	// Get retrieves a template by name
	Get(name string) (*Template, error)

	// List returns all available templates
	List() ([]*Template, error)

	// ListByCategory returns templates in a specific category
	ListByCategory(category string) ([]*Template, error)

	// Exists checks if a template exists
	Exists(name string) bool

	// Save saves a template (for writable repositories)
	Save(template *Template) error

	// Delete deletes a template (for writable repositories)
	Delete(name string) error
}

// TemplateNotFoundError is returned when a template is not found
type TemplateNotFoundError struct {
	Name string
}

func (e *TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template not found: %s", e.Name)
}

// TemplateAlreadyExistsError is returned when trying to save a template that already exists
type TemplateAlreadyExistsError struct {
	Name string
}

func (e *TemplateAlreadyExistsError) Error() string {
	return fmt.Sprintf("template already exists: %s", e.Name)
}

// ReadOnlyRepositoryError is returned when trying to modify a read-only repository
type ReadOnlyRepositoryError struct {
	Operation string
}

func (e *ReadOnlyRepositoryError) Error() string {
	return fmt.Sprintf("repository is read-only, cannot perform: %s", e.Operation)
}
