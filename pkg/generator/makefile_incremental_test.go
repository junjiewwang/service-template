package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	configtestutil "github.com/junjiewwang/service-template/pkg/config/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakefileIncrementalUpdate tests that Makefile generation uses incremental update strategy
func TestMakefileIncrementalUpdate(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test config
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	// Create generator
	gen := NewGenerator(cfg, tmpDir)

	// Step 1: Write initial user content to Makefile
	makefilePath := filepath.Join(tmpDir, "Makefile")
	userContent := `# User-defined targets
.PHONY: custom-target
custom-target:
	@echo "This is a custom target"

.PHONY: another-custom
another-custom:
	@echo "Another custom target"
`
	err := os.WriteFile(makefilePath, []byte(userContent), 0644)
	require.NoError(t, err, "Failed to write initial user content")

	// Step 2: Generate Makefile (should append generated content with markers)
	err = gen.generateMakefile()
	require.NoError(t, err, "Failed to generate Makefile")

	// Step 3: Read the generated Makefile
	content, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "Failed to read generated Makefile")
	contentStr := string(content)

	// Step 4: Verify that user content is preserved
	assert.Contains(t, contentStr, "# User-defined targets", "User content should be preserved")
	assert.Contains(t, contentStr, "custom-target", "User target should be preserved")
	assert.Contains(t, contentStr, "another-custom", "User target should be preserved")

	// Step 5: Verify that generated content has markers
	assert.Contains(t, contentStr, "# ===== GENERATED_START =====", "Start marker should be present")
	assert.Contains(t, contentStr, "# ===== GENERATED_END =====", "End marker should be present")

	// Step 6: Verify that generated content is present
	// The generated Makefile should contain common targets like build, test, etc.
	assert.Contains(t, contentStr, ".PHONY:", "Generated content should contain targets")

	// Step 7: Save the first generation result
	firstGeneration := contentStr

	// Step 8: Generate again (should be idempotent)
	err = gen.generateMakefile()
	require.NoError(t, err, "Failed to generate Makefile second time")

	// Step 9: Read the Makefile again
	content, err = os.ReadFile(makefilePath)
	require.NoError(t, err, "Failed to read Makefile after second generation")
	secondGeneration := string(content)

	// Step 10: Verify idempotency (content should be identical)
	assert.Equal(t, firstGeneration, secondGeneration, "Multiple generations should be idempotent")

	// Step 11: Verify structure
	// The file should have: user content + marker start + generated content + marker end
	lines := strings.Split(contentStr, "\n")

	var hasUserContent bool
	var hasStartMarker bool
	var hasEndMarker bool
	var hasGeneratedContent bool

	for _, line := range lines {
		if strings.Contains(line, "# User-defined targets") {
			hasUserContent = true
		}
		if strings.Contains(line, "# ===== GENERATED_START =====") {
			hasStartMarker = true
		}
		if strings.Contains(line, "# ===== GENERATED_END =====") {
			hasEndMarker = true
		}
		if hasStartMarker && !hasEndMarker && strings.Contains(line, ".PHONY:") {
			hasGeneratedContent = true
		}
	}

	assert.True(t, hasUserContent, "Should have user content")
	assert.True(t, hasStartMarker, "Should have start marker")
	assert.True(t, hasEndMarker, "Should have end marker")
	assert.True(t, hasGeneratedContent, "Should have generated content between markers")

	t.Logf("Generated Makefile structure verified successfully")
	t.Logf("Total lines: %d", len(lines))
}

// TestMakefileIncrementalUpdateWithoutExisting tests Makefile generation when no existing file
func TestMakefileIncrementalUpdateWithoutExisting(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test config
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	// Create generator
	gen := NewGenerator(cfg, tmpDir)

	// Generate Makefile (no existing file)
	err := gen.generateMakefile()
	require.NoError(t, err, "Failed to generate Makefile")

	// Read the generated Makefile
	makefilePath := filepath.Join(tmpDir, "Makefile")
	content, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "Failed to read generated Makefile")
	contentStr := string(content)

	// Even when no existing file, the content should have markers
	// This ensures consistency and allows proper incremental updates
	assert.NotEmpty(t, contentStr, "Generated Makefile should not be empty")
	assert.Contains(t, contentStr, ".PHONY:", "Generated Makefile should contain targets")
	assert.Contains(t, contentStr, "# ===== GENERATED_START =====", "First generation should have start marker")
	assert.Contains(t, contentStr, "# ===== GENERATED_END =====", "First generation should have end marker")

	t.Logf("Generated Makefile without existing file: %d bytes", len(content))

	// Step 2: Generate again to verify idempotency
	err = gen.generateMakefile()
	require.NoError(t, err, "Failed to generate Makefile second time")

	// Read again
	content2, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "Failed to read Makefile after second generation")

	// Should be identical (idempotent)
	assert.Equal(t, string(content), string(content2), "Multiple generations should be idempotent")

	t.Logf("Verified idempotency: content unchanged after second generation")
}

// TestMakefileIncrementalUpdatePreservesUserChanges tests that user changes outside markers are preserved
func TestMakefileIncrementalUpdatePreservesUserChanges(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test config
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	// Create generator
	gen := NewGenerator(cfg, tmpDir)

	makefilePath := filepath.Join(tmpDir, "Makefile")

	// Step 1: Write initial Makefile with user content
	initialContent := `# My custom Makefile header
# Author: Test User

.PHONY: my-custom-target
my-custom-target:
	@echo "Custom target before generation"
`
	err := os.WriteFile(makefilePath, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Step 2: First generation
	err = gen.generateMakefile()
	require.NoError(t, err)

	// Step 3: Read and verify
	content, err := os.ReadFile(makefilePath)
	require.NoError(t, err)
	contentStr := string(content)

	// Verify user content is preserved
	assert.Contains(t, contentStr, "# My custom Makefile header")
	assert.Contains(t, contentStr, "# Author: Test User")
	assert.Contains(t, contentStr, "my-custom-target")

	// Step 4: Add more user content after the generated block
	additionalUserContent := `
# Additional user content added after generation
.PHONY: post-generation-target
post-generation-target:
	@echo "This was added after generation"
`
	updatedContent := contentStr + additionalUserContent
	err = os.WriteFile(makefilePath, []byte(updatedContent), 0644)
	require.NoError(t, err)

	// Step 5: Generate again
	err = gen.generateMakefile()
	require.NoError(t, err)

	// Step 6: Verify all user content is still preserved
	finalContent, err := os.ReadFile(makefilePath)
	require.NoError(t, err)
	finalContentStr := string(finalContent)

	assert.Contains(t, finalContentStr, "# My custom Makefile header", "Original user header should be preserved")
	assert.Contains(t, finalContentStr, "my-custom-target", "Original user target should be preserved")
	assert.Contains(t, finalContentStr, "post-generation-target", "Additional user target should be preserved")
	assert.Contains(t, finalContentStr, "# ===== GENERATED_START =====", "Generated block should still exist")
	assert.Contains(t, finalContentStr, "# ===== GENERATED_END =====", "Generated block should still exist")

	t.Logf("User changes preserved successfully across multiple generations")
	t.Logf("Final Makefile: %d lines", len(strings.Split(finalContentStr, "\n")))
}
