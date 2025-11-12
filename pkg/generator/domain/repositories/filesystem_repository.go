package repositories

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileSystemTemplateRepository stores templates in the file system
type FileSystemTemplateRepository struct {
	basePath  string
	templates map[string]*Template
	mu        sync.RWMutex
	loaded    bool
}

// NewFileSystemTemplateRepository creates a new file system template repository
func NewFileSystemTemplateRepository(basePath string) (*FileSystemTemplateRepository, error) {
	// Validate base path
	if basePath == "" {
		return nil, fmt.Errorf("base path cannot be empty")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	repo := &FileSystemTemplateRepository{
		basePath:  basePath,
		templates: make(map[string]*Template),
	}

	// Load templates from file system
	if err := repo.load(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return repo, nil
}

// load loads all templates from the file system
func (r *FileSystemTemplateRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := filepath.Walk(r.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Read template content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Extract template name (relative to base path, without .tmpl extension)
		relPath, err := filepath.Rel(r.basePath, path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(relPath, ".tmpl")

		// Determine category from directory structure
		category := "default"
		if dir := filepath.Dir(relPath); dir != "." {
			category = dir
		}

		template := &Template{
			Name:     name,
			Content:  string(content),
			Category: category,
		}

		r.templates[name] = template
		return nil
	})

	if err != nil {
		return err
	}

	r.loaded = true
	return nil
}

// Reload reloads all templates from the file system
func (r *FileSystemTemplateRepository) Reload() error {
	r.mu.Lock()
	r.templates = make(map[string]*Template)
	r.mu.Unlock()

	return r.load()
}

// Get retrieves a template by name
func (r *FileSystemTemplateRepository) Get(name string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.templates[name]
	if !exists {
		return nil, &TemplateNotFoundError{Name: name}
	}

	return template, nil
}

// List returns all available templates
func (r *FileSystemTemplateRepository) List() ([]*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]*Template, 0, len(r.templates))
	for _, template := range r.templates {
		templates = append(templates, template)
	}

	return templates, nil
}

// ListByCategory returns templates in a specific category
func (r *FileSystemTemplateRepository) ListByCategory(category string) ([]*Template, error) {
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
func (r *FileSystemTemplateRepository) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.templates[name]
	return exists
}

// Save saves a template to the file system
func (r *FileSystemTemplateRepository) Save(template *Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Construct file path
	filePath := filepath.Join(r.basePath, template.Name+".tmpl")

	// Create directory if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write template content
	if err := os.WriteFile(filePath, []byte(template.Content), 0644); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	// Update in-memory cache
	r.templates[template.Name] = template

	return nil
}

// Delete deletes a template from the file system
func (r *FileSystemTemplateRepository) Delete(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if template exists
	if _, exists := r.templates[name]; !exists {
		return &TemplateNotFoundError{Name: name}
	}

	// Construct file path
	filePath := filepath.Join(r.basePath, name+".tmpl")

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	// Remove from in-memory cache
	delete(r.templates, name)

	return nil
}

// GetBasePath returns the base path
func (r *FileSystemTemplateRepository) GetBasePath() string {
	return r.basePath
}

// Count returns the number of templates
func (r *FileSystemTemplateRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.templates)
}
