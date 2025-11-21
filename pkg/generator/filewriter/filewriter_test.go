package filewriter_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/filewriter"
	"github.com/junjiewwang/service-template/pkg/generator/filewriter/strategies"
)

func TestFileWriter_Overwrite(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	ctx := context.Background()

	// Write initial content
	writer := filewriter.New()
	err := writer.WriteString(ctx, testFile, "initial content")
	if err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	// Verify initial content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "initial content" {
		t.Errorf("Expected 'initial content', got '%s'", string(content))
	}

	// Overwrite with new content
	err = writer.WriteString(ctx, testFile, "new content")
	if err != nil {
		t.Fatalf("Failed to overwrite content: %v", err)
	}

	// Verify overwritten content
	content, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "new content" {
		t.Errorf("Expected 'new content', got '%s'", string(content))
	}
}

func TestFileWriter_Skip(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	ctx := context.Background()

	// Write initial content with skip strategy
	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.SkipStrategyID))

	err := writer.WriteString(ctx, testFile, "initial content")
	if err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	// Try to write again (should be skipped)
	err = writer.WriteString(ctx, testFile, "new content")
	if err != nil {
		t.Fatalf("Failed to skip write: %v", err)
	}

	// Verify content is still the initial content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "initial content" {
		t.Errorf("Expected 'initial content', got '%s'", string(content))
	}
}

func TestFileWriter_Incremental(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	ctx := context.Background()

	// Write initial content
	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.IncrementalStrategyID))

	initialContent := "# User content\nuser_var = 1\n"
	err := writer.WriteString(ctx, testFile, initialContent)
	if err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	// Write generated content (should append with markers)
	generatedContent := "generated_var = 2"
	err = writer.WriteString(ctx, testFile, generatedContent)
	if err != nil {
		t.Fatalf("Failed to write generated content: %v", err)
	}

	// Verify content has markers
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "# User content") {
		t.Error("User content should be preserved")
	}
	if !contains(contentStr, "# ===== GENERATED_START =====") {
		t.Error("Start marker should be present")
	}
	if !contains(contentStr, "generated_var = 2") {
		t.Error("Generated content should be present")
	}
	if !contains(contentStr, "# ===== GENERATED_END =====") {
		t.Error("End marker should be present")
	}

	// Write again with updated content (should replace marker block)
	updatedContent := "generated_var = 3"
	err = writer.WriteString(ctx, testFile, updatedContent)
	if err != nil {
		t.Fatalf("Failed to write updated content: %v", err)
	}

	// Verify content is updated
	content, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr = string(content)
	t.Logf("Final content:\n%s", contentStr)
	if !contains(contentStr, "# User content") {
		t.Error("User content should still be preserved")
	}
	if !contains(contentStr, "generated_var = 3") {
		t.Error("Generated content should be updated")
	}
	if contains(contentStr, "generated_var = 2") {
		t.Error("Old generated content should be replaced")
	}
}

func TestFileWriter_IncrementalIdempotent(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	ctx := context.Background()

	// Write with incremental strategy
	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.IncrementalStrategyID))

	// Write initial content with existing user content
	initialContent := "# User content\nuser_var = 1\n"
	err := os.WriteFile(testFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	content := "test content"

	// Write first time (should append with markers)
	err = writer.WriteString(ctx, testFile, content)
	if err != nil {
		t.Fatalf("Failed to write first time: %v", err)
	}

	firstContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Write second time with same content (should be idempotent)
	err = writer.WriteString(ctx, testFile, content)
	if err != nil {
		t.Fatalf("Failed to write second time: %v", err)
	}

	secondContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Content should be identical (idempotent)
	if string(firstContent) != string(secondContent) {
		t.Errorf("Content should be identical after second write (idempotent)")
		t.Logf("First content:\n%s", string(firstContent))
		t.Logf("Second content:\n%s", string(secondContent))
	}

	// Write third time with different content (should update)
	newContent := "updated content"
	err = writer.WriteString(ctx, testFile, newContent)
	if err != nil {
		t.Fatalf("Failed to write third time: %v", err)
	}

	thirdContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Content should be different from second write
	if string(secondContent) == string(thirdContent) {
		t.Error("Content should be updated after third write")
	}

	// Should contain updated content
	if !contains(string(thirdContent), "updated content") {
		t.Error("Should contain updated content")
	}

	// Should still preserve user content
	if !contains(string(thirdContent), "# User content") {
		t.Error("Should still preserve user content")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
