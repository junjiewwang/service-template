package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator"
)

// TestEndToEndGeneration tests the complete generation flow
func TestEndToEndGeneration(t *testing.T) {
	// Create temporary output directory
	tmpDir, err := os.MkdirTemp("", "svcgen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test configuration
	cfg := createTestConfig()

	// Create generator
	gen := generator.NewGenerator(cfg, tmpDir)

	// Execute generation
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify generated files
	serviceName := cfg.Service.Name
	expectedFiles := []string{
		"Makefile",
		"compose.yaml",
		".tad/devops.yaml",
		filepath.Join(".tad/build", serviceName, "build.sh"),
		filepath.Join(".tad/build", serviceName, "build_deps_install.sh"),
		filepath.Join(".tad/build", serviceName, "entrypoint.sh"),
		filepath.Join(".tad/build", serviceName, "rt_prepare.sh"),
	}

	for _, file := range expectedFiles {
		path := filepath.Join(tmpDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file not generated: %s", file)
		} else {
			t.Logf("✓ Generated: %s", file)
		}
	}

	// Verify Dockerfiles for each architecture
	for _, arch := range []string{"amd64", "arm64"} {
		dockerfilePath := filepath.Join(tmpDir, ".tad", "build", serviceName,
			"Dockerfile."+serviceName+"."+arch)
		if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
			t.Errorf("Expected Dockerfile not generated: %s", dockerfilePath)
		} else {
			t.Logf("✓ Generated: Dockerfile.%s.%s", serviceName, arch)
		}
	}

	// Verify file permissions for scripts
	scripts := []string{
		filepath.Join(".tad/build", serviceName, "build.sh"),
		filepath.Join(".tad/build", serviceName, "build_deps_install.sh"),
		filepath.Join(".tad/build", serviceName, "entrypoint.sh"),
		filepath.Join(".tad/build", serviceName, "rt_prepare.sh"),
	}

	for _, script := range scripts {
		path := filepath.Join(tmpDir, script)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		mode := info.Mode()
		if mode&0111 == 0 {
			t.Errorf("Script not executable: %s (mode: %v)", script, mode)
		} else {
			t.Logf("✓ Executable: %s", script)
		}
	}
}

// TestGenerationWithPlugins tests generation with plugins configured
func TestGenerationWithPlugins(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "svcgen-test-plugins-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := createTestConfig()

	// Add plugins
	cfg.Plugins.InstallDir = "/opt/plugins"
	cfg.Plugins.Items = []config.PluginConfig{
		{
			Name:           "test-plugin",
			Description:    "Test plugin",
			DownloadURL:    config.NewStaticDownloadURL("https://example.com/plugin.tar.gz"),
			InstallCommand: "tar -xzf plugin.tar.gz",
			Required:       true,
		},
	}

	gen := generator.NewGenerator(cfg, tmpDir)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify build_plugins.sh is generated
	serviceName := cfg.Service.Name
	pluginScript := filepath.Join(tmpDir, ".tad", "build", serviceName, "build_plugins.sh")
	if _, err := os.Stat(pluginScript); os.IsNotExist(err) {
		t.Errorf("build_plugins.sh not generated")
	} else {
		t.Logf("✓ Generated: build_plugins.sh")
	}
}

// TestGenerationWithHealthcheck tests generation with healthcheck enabled
func TestGenerationWithHealthcheck(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "svcgen-test-healthcheck-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := createTestConfig()
	cfg.Runtime.Healthcheck.Enabled = true
	cfg.Runtime.Healthcheck.Type = "default"

	gen := generator.NewGenerator(cfg, tmpDir)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify healthcheck.sh is generated
	serviceName := cfg.Service.Name
	healthcheckScript := filepath.Join(tmpDir, ".tad", "build", serviceName, "healthchk.sh")
	if _, err := os.Stat(healthcheckScript); os.IsNotExist(err) {
		t.Errorf("healthcheck.sh not generated")
	} else {
		t.Logf("✓ Generated: healthcheck.sh")
	}
}

// TestGenerationWithCustomHealthcheck tests generation with custom healthcheck
func TestGenerationWithCustomHealthcheck(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "svcgen-test-custom-healthcheck-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := createTestConfig()
	cfg.Runtime.Healthcheck.Enabled = true
	cfg.Runtime.Healthcheck.Type = "custom"
	cfg.Runtime.Healthcheck.CustomScript = "curl -f http://localhost:8080/health || exit 1"

	gen := generator.NewGenerator(cfg, tmpDir)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify healthcheck.sh is generated with custom script
	serviceName := cfg.Service.Name
	healthcheckScript := filepath.Join(tmpDir, ".tad", "build", serviceName, "healthchk.sh")
	content, err := os.ReadFile(healthcheckScript)
	if err != nil {
		t.Fatalf("Failed to read healthcheck.sh: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("healthcheck.sh is empty")
	} else {
		t.Logf("✓ Generated custom healthcheck.sh (%d bytes)", len(content))
	}
}

// TestGenerationWithDifferentLanguages tests generation for different languages
func TestGenerationWithDifferentLanguages(t *testing.T) {
	languages := []struct {
		name     string
		langType string
		config   map[string]interface{}
	}{
		{
			name:     "Go",
			langType: "go",
			config: map[string]interface{}{
				"version": "1.21",
				"module":  "github.com/example/test-service",
			},
		},
		{
			name:     "Python",
			langType: "python",
			config: map[string]interface{}{
				"version":           "3.11",
				"requirements_file": "requirements.txt",
			},
		},
		{
			name:     "Node.js",
			langType: "nodejs",
			config: map[string]interface{}{
				"version":      "20",
				"package_file": "package.json",
			},
		},
	}

	for _, lang := range languages {
		t.Run(lang.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "svcgen-test-"+lang.langType+"-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			cfg := createTestConfig()
			cfg.Language.Type = lang.langType
			cfg.Language.Config = lang.config

			gen := generator.NewGenerator(cfg, tmpDir)
			if err := gen.Generate(); err != nil {
				t.Fatalf("Generation failed for %s: %v", lang.name, err)
			}

			t.Logf("✓ Successfully generated for %s", lang.name)
		})
	}
}

// TestGenerationFileContent tests the content of generated files
func TestGenerationFileContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "svcgen-test-content-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := createTestConfig()
	gen := generator.NewGenerator(cfg, tmpDir)

	if err := gen.Generate(); err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Test Makefile content
	makefilePath := filepath.Join(tmpDir, "Makefile")
	makefileContent, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}

	if len(makefileContent) == 0 {
		t.Errorf("Makefile is empty")
	} else {
		t.Logf("✓ Makefile generated (%d bytes)", len(makefileContent))
	}

	// Test compose.yaml content
	composePath := filepath.Join(tmpDir, "compose.yaml")
	composeContent, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("Failed to read compose.yaml: %v", err)
	}

	if len(composeContent) == 0 {
		t.Errorf("compose.yaml is empty")
	} else {
		t.Logf("✓ compose.yaml generated (%d bytes)", len(composeContent))
	}
}

// createTestConfig creates a test configuration
func createTestConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:        "test-service",
			Description: "Test service for integration testing",
			Ports: []config.PortConfig{
				{
					Name:     "http",
					Port:     8080,
					Protocol: "TCP",
					Expose:   true,
				},
			},
			DeployDir: "/opt/services",
		},
		Language: config.LanguageConfig{
			Type: "go",
			Config: map[string]interface{}{
				"version": "1.21",
				"module":  "github.com/example/test-service",
			},
		},
		Build: config.BuildConfig{
			DependencyFiles: config.DependencyFilesConfig{
				AutoDetect: true,
				Files:      []string{"go.mod", "go.sum"},
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
				ARM64: "golang:1.21-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:3.18",
				ARM64: "alpine:3.18",
			},
			Dependencies: config.DependenciesConfig{
				SystemPkgs: []string{"ca-certificates", "tzdata"},
			},
			Commands: config.BuildCommandsConfig{
				Build: "go build -o bin/service ./cmd/service",
			},
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/opt/plugins",
			Items:      []config.PluginConfig{},
		},
		Runtime: config.RuntimeConfig{
			SystemDependencies: config.RuntimeSystemDependenciesConfig{
				Packages: []string{"ca-certificates"},
			},
			Healthcheck: config.HealthcheckConfig{
				Enabled: false,
				Type:    "default",
			},
			Startup: config.StartupConfig{
				Command: "/opt/services/test-service/bin/service",
			},
			GenerateScripts: true,
		},
		LocalDev: config.LocalDevConfig{
			Compose: config.ComposeConfig{
				Resources: config.ResourcesConfig{
					Limits: config.ResourceLimits{
						CPUs:   "2",
						Memory: "1G",
					},
				},
			},
		},
		Metadata: config.MetadataConfig{
			TemplateVersion: "1.0.0",
			Generator:       "svcgen",
		},
	}
}
