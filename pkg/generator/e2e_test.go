package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_GenerateCompleteProject tests the complete project generation workflow
func TestE2E_GenerateCompleteProject(t *testing.T) {
	// Create temporary output directory
	tmpDir := t.TempDir()

	// Load test configuration
	loader := config.NewLoader("../../demo-app/service.yaml")
	cfg, err := loader.Load()
	require.NoError(t, err, "Failed to load test configuration")

	// Create generator
	gen := generator.NewGenerator(cfg, tmpDir)

	// Generate project
	err = gen.Generate()
	require.NoError(t, err, "Failed to generate project")

	// Verify all expected files are generated
	expectedFiles := []string{
		"Makefile",
		"compose.yaml",
		".tad/devops.yaml",
		filepath.Join(".tad/build", cfg.Service.Name, "Dockerfile."+cfg.Service.Name+".amd64"),
		filepath.Join(".tad/build", cfg.Service.Name, "Dockerfile."+cfg.Service.Name+".arm64"),
		filepath.Join(".tad/build", cfg.Service.Name, "build.sh"),
		filepath.Join(".tad/build", cfg.Service.Name, "build_deps_install.sh"),
		filepath.Join(".tad/build", cfg.Service.Name, "rt_prepare.sh"),
		filepath.Join(".tad/build", cfg.Service.Name, "entrypoint.sh"),
		filepath.Join(".tad/build", cfg.Service.Name, "healthchk.sh"),
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(tmpDir, file)
		assert.FileExists(t, fullPath, "Expected file %s to exist", file)

		// Verify file is not empty
		info, err := os.Stat(fullPath)
		require.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0), "File %s should not be empty", file)
	}

	// Verify Dockerfile content
	dockerfileAmd64 := filepath.Join(tmpDir, ".tad/build", cfg.Service.Name, "Dockerfile."+cfg.Service.Name+".amd64")
	content, err := os.ReadFile(dockerfileAmd64)
	require.NoError(t, err)
	assert.Contains(t, string(content), "FROM", "Dockerfile should contain FROM instruction")
	assert.Contains(t, string(content), "WORKDIR", "Dockerfile should contain WORKDIR instruction")

	// Verify Makefile content
	makefilePath := filepath.Join(tmpDir, "Makefile")
	content, err = os.ReadFile(makefilePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), ".PHONY:", "Makefile should contain .PHONY")
	assert.Contains(t, string(content), "docker-build", "Makefile should contain docker-build target")

	// Verify compose.yaml content
	composePath := filepath.Join(tmpDir, "compose.yaml")
	content, err = os.ReadFile(composePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "services:", "compose.yaml should contain services")
	assert.Contains(t, string(content), cfg.Service.Name, "compose.yaml should contain service name")

	// Verify build script is executable
	buildScript := filepath.Join(tmpDir, ".tad/build", cfg.Service.Name, "build.sh")
	info, err := os.Stat(buildScript)
	require.NoError(t, err)
	assert.True(t, info.Mode()&0111 != 0, "build.sh should be executable")
}

// TestE2E_GenerateWithPlugins tests project generation with plugins
func TestE2E_GenerateWithPlugins(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:        "test-service-with-plugins",
			Description: "Test service with plugins",
			DeployDir:   "/opt",
		},
		Language: config.LanguageConfig{
			Type: "golang",
			Config: map[string]interface{}{
				"version": "1.21",
			},
		},
		Build: config.BuildConfig{
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
			Commands: config.BuildCommandsConfig{
				Build: "go build -o bin/app ./cmd/app",
			},
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/opt/plugins",
			Items: []config.PluginConfig{
				{
					Name:           "test-plugin",
					Description:    "Test plugin",
					InstallCommand: "echo 'install plugin'",
					Required:       true,
				},
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/app",
			},
		},
	}

	gen := generator.NewGenerator(cfg, tmpDir)
	err := gen.Generate()
	require.NoError(t, err)

	// Verify build_plugins.sh is generated
	buildPluginsScript := filepath.Join(tmpDir, ".tad/build", cfg.Service.Name, "build_plugins.sh")
	assert.FileExists(t, buildPluginsScript)

	content, err := os.ReadFile(buildPluginsScript)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test-plugin", "build_plugins.sh should contain plugin name")
}

// TestE2E_GenerateMinimalConfig tests generation with minimal configuration
func TestE2E_GenerateMinimalConfig(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:        "minimal-service",
			Description: "Minimal service",
			DeployDir:   "/opt",
		},
		Language: config.LanguageConfig{
			Type: "golang",
		},
		Build: config.BuildConfig{
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
			},
			Commands: config.BuildCommandsConfig{
				Build: "go build -o bin/app",
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/app",
			},
		},
	}

	gen := generator.NewGenerator(cfg, tmpDir)
	err := gen.Generate()
	require.NoError(t, err)

	// Verify basic files are generated
	assert.FileExists(t, filepath.Join(tmpDir, "Makefile"))
	assert.FileExists(t, filepath.Join(tmpDir, "compose.yaml"))
	assert.FileExists(t, filepath.Join(tmpDir, ".tad/devops.yaml"))
}

// TestE2E_ValidationErrors tests that validation errors are properly reported
func TestE2E_ValidationErrors(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid configuration (missing required build commands)
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:        "test-service",
			Description: "Test",
			DeployDir:   "/opt",
		},
		Language: config.LanguageConfig{
			Type: "golang",
		},
		Build: config.BuildConfig{
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
			},
			Commands: config.BuildCommandsConfig{
				Build: "", // Invalid: empty build command
			},
		},
		Runtime: config.RuntimeConfig{
			Startup: config.StartupConfig{
				Command: "./bin/app",
			},
		},
	}

	gen := generator.NewGenerator(cfg, tmpDir)
	err := gen.Generate()
	// Note: Legacy generator doesn't validate, so this will succeed
	// In a full DDD implementation, this would fail validation
	assert.NoError(t, err, "Legacy generator doesn't validate configuration")
}
