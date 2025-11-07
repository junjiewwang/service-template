package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestIntegration tests the complete workflow
func TestIntegration(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "svcgen"), "./cmd/svcgen")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	binary := filepath.Join(tmpDir, "svcgen")
	testDir := filepath.Join(tmpDir, "test-project")

	// Create test directory
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test init command
	t.Run("init", func(t *testing.T) {
		cmd := exec.Command(binary, "init", "-c", filepath.Join(testDir, "service.yaml"))
		cmd.Dir = testDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Init command failed: %v", err)
		}

		// Check if service.yaml was created
		if _, err := os.Stat(filepath.Join(testDir, "service.yaml")); os.IsNotExist(err) {
			t.Fatal("service.yaml was not created")
		}
	})

	// Copy example config
	exampleConfig := "service.yaml.example"
	if data, err := os.ReadFile(exampleConfig); err == nil {
		if err := os.WriteFile(filepath.Join(testDir, "service.yaml"), data, 0644); err != nil {
			t.Fatalf("Failed to copy example config: %v", err)
		}
	}

	// Test validate command
	t.Run("validate", func(t *testing.T) {
		cmd := exec.Command(binary, "validate", "-c", filepath.Join(testDir, "service.yaml"))
		cmd.Dir = testDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Validate command failed: %v", err)
		}
	})

	// Test generate command
	t.Run("generate", func(t *testing.T) {
		cmd := exec.Command(binary, "generate", "-c", filepath.Join(testDir, "service.yaml"), "-o", testDir)
		cmd.Dir = testDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Generate command failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("Generate output:\n%s", string(output))

		// List all generated files for debugging
		filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				relPath, _ := filepath.Rel(testDir, path)
				t.Logf("Generated file: %s", relPath)
			}
			return nil
		})

		// Check if expected files were generated
		expectedFiles := []string{
			".tad/build/example-service/Dockerfile.example-service.amd64",
			".tad/build/example-service/Dockerfile.example-service.arm64",
			".tad/build/example-service/build.sh",
			".tad/build/example-service/build_deps_install.sh",
			".tad/build/example-service/rt_prepare.sh",
			".tad/build/example-service/entrypoint.sh",
			".tad/build/example-service/healthchk.sh",
			"compose.yaml",
			"Makefile",
			".tad/devops.yaml",
		}

		for _, file := range expectedFiles {
			path := filepath.Join(testDir, file)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not generated", file)
			}
		}
	})
}
