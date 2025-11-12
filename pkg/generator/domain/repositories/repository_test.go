package repositories

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedTemplateRepository(t *testing.T) {
	repo := NewEmbeddedTemplateRepository()

	// Test Register
	template := &Template{
		Name:        "test-template",
		Content:     "Hello {{.Name}}",
		Description: "Test template",
		Category:    "test",
		Tags:        []string{"test"},
	}

	err := repo.Register(template)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Test duplicate registration
	err = repo.Register(template)
	if _, ok := err.(*TemplateAlreadyExistsError); !ok {
		t.Error("Register() should return TemplateAlreadyExistsError for duplicate")
	}

	// Test Get
	retrieved, err := repo.Get("test-template")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if retrieved.Name != template.Name {
		t.Errorf("Get() returned wrong template")
	}

	// Test Get non-existent
	_, err = repo.Get("non-existent")
	if _, ok := err.(*TemplateNotFoundError); !ok {
		t.Error("Get() should return TemplateNotFoundError for non-existent template")
	}

	// Test Exists
	if !repo.Exists("test-template") {
		t.Error("Exists() should return true for existing template")
	}

	if repo.Exists("non-existent") {
		t.Error("Exists() should return false for non-existent template")
	}

	// Test List
	templates, err := repo.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("List() returned %v templates, want 1", len(templates))
	}

	// Test ListByCategory
	templates, err = repo.ListByCategory("test")
	if err != nil {
		t.Fatalf("ListByCategory() error = %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("ListByCategory() returned %v templates, want 1", len(templates))
	}

	templates, err = repo.ListByCategory("other")
	if err != nil {
		t.Fatalf("ListByCategory() error = %v", err)
	}

	if len(templates) != 0 {
		t.Errorf("ListByCategory() returned %v templates, want 0", len(templates))
	}

	// Test Count
	if repo.Count() != 1 {
		t.Errorf("Count() = %v, want 1", repo.Count())
	}

	// Test Save (should fail - read-only)
	err = repo.Save(template)
	if _, ok := err.(*ReadOnlyRepositoryError); !ok {
		t.Error("Save() should return ReadOnlyRepositoryError")
	}

	// Test Delete (should fail - read-only)
	err = repo.Delete("test-template")
	if _, ok := err.(*ReadOnlyRepositoryError); !ok {
		t.Error("Delete() should return ReadOnlyRepositoryError")
	}

	// Test Clear
	repo.Clear()
	if repo.Count() != 0 {
		t.Error("Clear() did not remove all templates")
	}
}

func TestFileSystemTemplateRepository(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "template-repo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template file
	templateContent := "Hello {{.Name}}"
	templatePath := filepath.Join(tmpDir, "test-template.tmpl")
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test template: %v", err)
	}

	// Create repository
	repo, err := NewFileSystemTemplateRepository(tmpDir)
	if err != nil {
		t.Fatalf("NewFileSystemTemplateRepository() error = %v", err)
	}

	// Test Get
	template, err := repo.Get("test-template")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if template.Content != templateContent {
		t.Errorf("Get() returned wrong content: %v, want %v", template.Content, templateContent)
	}

	// Test Exists
	if !repo.Exists("test-template") {
		t.Error("Exists() should return true")
	}

	// Test List
	templates, err := repo.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("List() returned %v templates, want 1", len(templates))
	}

	// Test Save
	newTemplate := &Template{
		Name:     "new-template",
		Content:  "New content",
		Category: "test",
	}

	err = repo.Save(newTemplate)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	newPath := filepath.Join(tmpDir, "new-template.tmpl")
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Error("Save() did not create file")
	}

	// Test Get after Save
	retrieved, err := repo.Get("new-template")
	if err != nil {
		t.Fatalf("Get() after Save() error = %v", err)
	}

	if retrieved.Content != newTemplate.Content {
		t.Error("Get() after Save() returned wrong content")
	}

	// Test Count
	if repo.Count() != 2 {
		t.Errorf("Count() = %v, want 2", repo.Count())
	}

	// Test Delete
	err = repo.Delete("new-template")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Error("Delete() did not remove file")
	}

	// Test Get after Delete
	_, err = repo.Get("new-template")
	if _, ok := err.(*TemplateNotFoundError); !ok {
		t.Error("Get() after Delete() should return TemplateNotFoundError")
	}

	// Test Reload
	err = repo.Reload()
	if err != nil {
		t.Fatalf("Reload() error = %v", err)
	}

	if repo.Count() != 1 {
		t.Errorf("Count() after Reload() = %v, want 1", repo.Count())
	}
}

func TestFileSystemTemplateRepository_Categories(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "template-repo-cat-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create templates in different categories
	categories := []string{"docker", "scripts", "config"}
	for _, cat := range categories {
		catDir := filepath.Join(tmpDir, cat)
		os.MkdirAll(catDir, 0755)

		templatePath := filepath.Join(catDir, "template.tmpl")
		os.WriteFile(templatePath, []byte("content"), 0644)
	}

	// Create repository
	repo, err := NewFileSystemTemplateRepository(tmpDir)
	if err != nil {
		t.Fatalf("NewFileSystemTemplateRepository() error = %v", err)
	}

	// Test ListByCategory
	for _, cat := range categories {
		templates, err := repo.ListByCategory(cat)
		if err != nil {
			t.Fatalf("ListByCategory(%s) error = %v", cat, err)
		}

		if len(templates) != 1 {
			t.Errorf("ListByCategory(%s) returned %v templates, want 1", cat, len(templates))
		}
	}

	// Test total count
	if repo.Count() != 3 {
		t.Errorf("Count() = %v, want 3", repo.Count())
	}
}

func TestFileSystemTemplateRepository_EmptyPath(t *testing.T) {
	_, err := NewFileSystemTemplateRepository("")
	if err == nil {
		t.Error("NewFileSystemTemplateRepository(\"\") should return error")
	}
}

func TestTemplateErrors(t *testing.T) {
	// Test TemplateNotFoundError
	err := &TemplateNotFoundError{Name: "test"}
	expected := "template not found: test"
	if err.Error() != expected {
		t.Errorf("TemplateNotFoundError.Error() = %v, want %v", err.Error(), expected)
	}

	// Test TemplateAlreadyExistsError
	err2 := &TemplateAlreadyExistsError{Name: "test"}
	expected2 := "template already exists: test"
	if err2.Error() != expected2 {
		t.Errorf("TemplateAlreadyExistsError.Error() = %v, want %v", err2.Error(), expected2)
	}

	// Test ReadOnlyRepositoryError
	err3 := &ReadOnlyRepositoryError{Operation: "Save"}
	expected3 := "repository is read-only, cannot perform: Save"
	if err3.Error() != expected3 {
		t.Errorf("ReadOnlyRepositoryError.Error() = %v, want %v", err3.Error(), expected3)
	}
}
