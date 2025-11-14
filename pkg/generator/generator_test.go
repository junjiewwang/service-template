package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	// Arrange: Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	t.Logf("Created temp directory: %s", tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{},
		Language: config.LanguageConfig{
			Type: "golang",
		},
		LocalDev: config.LocalDevConfig{
			Kubernetes: config.KubernetesConfig{
				Enabled: true,
			},
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/opt/plugins",
			Items: []config.PluginConfig{
				{
					Name:        "test-plugin",
					Description: "Test plugin",
					DownloadURL: config.NewStaticDownloadURL("https://example.com/plugin.tar.gz"),
				},
			},
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Act: Generate all files
	err = gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")

	// Assert: Check that expected files were created
	ciPaths := context.NewCIPaths(cfg)
	expectedFiles := []string{
		".tad/build/test-service/Dockerfile.test-service.amd64",
		".tad/build/test-service/Dockerfile.test-service.arm64",
		"compose.yaml",
		"Makefile",
		ciPaths.GetScriptPath(ciPaths.BuildScript),
		ciPaths.GetScriptPath(ciPaths.DepsInstallScript),
		ciPaths.GetScriptPath(ciPaths.RtPrepareScript),
		ciPaths.GetScriptPath(ciPaths.EntrypointScript),
		ciPaths.GetScriptPath(ciPaths.HealthcheckScript),
		".tad/devops.yaml",
	}

	// List all generated files for debugging
	var generatedFiles []string
	filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			relPath, _ := filepath.Rel(outputDir, path)
			generatedFiles = append(generatedFiles, relPath)
		}
		return nil
	})
	t.Logf("Generated %d files in total", len(generatedFiles))

	// Verify each expected file exists
	for _, file := range expectedFiles {
		fullPath := filepath.Join(outputDir, file)
		_, err := os.Stat(fullPath)
		require.NoError(t, err, "Expected file should exist: %s", file)
		t.Logf("✓ Verified file exists: %s", file)
	}
}

func TestGenerator_GenerateDockerfiles(t *testing.T) {
	// Arrange: Setup test environment
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{},
		Language: config.LanguageConfig{
			Type: "golang",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Act: Generate Dockerfiles
	err = gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")

	// Assert: Check Dockerfiles were created
	dockerfileAMD64 := filepath.Join(outputDir, ".tad", "build", "test-service", "Dockerfile.test-service.amd64")
	dockerfileARM64 := filepath.Join(outputDir, ".tad", "build", "test-service", "Dockerfile.test-service.arm64")

	_, err = os.Stat(dockerfileAMD64)
	require.NoError(t, err, "Dockerfile.test-service.amd64 should be created")
	t.Logf("✓ Verified Dockerfile created: %s", dockerfileAMD64)

	_, err = os.Stat(dockerfileARM64)
	require.NoError(t, err, "Dockerfile.test-service.arm64 should be created")
	t.Logf("✓ Verified Dockerfile created: %s", dockerfileARM64)
}

func TestGenerator_GenerateCompose(t *testing.T) {
	// Arrange: Setup test environment
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{},
		Language: config.LanguageConfig{
			Type: "golang",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Act: Generate compose file
	err = gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")

	// Assert: Check compose file was created
	composePath := filepath.Join(outputDir, "compose.yaml")
	_, err = os.Stat(composePath)
	require.NoError(t, err, "compose.yaml should be created")
	t.Logf("✓ Verified compose.yaml created: %s", composePath)
}

func TestGenerator_GenerateMakefile(t *testing.T) {
	// Arrange: Setup test environment
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{},
		Language: config.LanguageConfig{
			Type: "golang",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Act: Generate Makefile
	err = gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")

	// Assert: Check Makefile was created
	makefilePath := filepath.Join(outputDir, "Makefile")
	_, err = os.Stat(makefilePath)
	require.NoError(t, err, "Makefile should be created")
	t.Logf("✓ Verified Makefile created: %s", makefilePath)
}

func TestGenerator_GenerateScripts(t *testing.T) {
	// Arrange: Setup test environment
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{},
		Language: config.LanguageConfig{
			Type: "golang",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Act: Generate scripts
	err = gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")

	// Assert: Check scripts were created
	ciPaths := context.NewCIPaths(cfg)
	expectedScripts := []string{
		ciPaths.GetScriptPath(ciPaths.BuildScript),
		ciPaths.GetScriptPath(ciPaths.DepsInstallScript),
		ciPaths.GetScriptPath(ciPaths.RtPrepareScript),
		ciPaths.GetScriptPath(ciPaths.EntrypointScript),
		ciPaths.GetScriptPath(ciPaths.HealthcheckScript),
	}

	t.Logf("Verifying %d scripts were created", len(expectedScripts))
	for _, script := range expectedScripts {
		scriptPath := filepath.Join(outputDir, script)
		_, err := os.Stat(scriptPath)
		require.NoError(t, err, "Script should be created: %s", script)
		t.Logf("✓ Verified script exists: %s", script)
	}
}

func TestGenerator_GenerateWithKubernetesConfig(t *testing.T) {
	// Arrange: Setup test environment with Kubernetes enabled
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{},
		Language: config.LanguageConfig{
			Type: "golang",
		},
		LocalDev: config.LocalDevConfig{
			Kubernetes: config.KubernetesConfig{
				Enabled: true,
			},
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Act: Generate project with Kubernetes configuration
	err = gen.Generate()
	require.NoError(t, err, "Generate() should not return an error")

	// Assert: Check Makefile was created with Kubernetes targets
	makefilePath := filepath.Join(outputDir, "Makefile")
	_, err = os.Stat(makefilePath)
	require.NoError(t, err, "Makefile should be created")

	// Verify Makefile contains Kubernetes-related targets
	makefileContent, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "Should be able to read Makefile")
	makefileStr := string(makefileContent)

	assert.Contains(t, makefileStr, "k8s-configmap",
		"Makefile should contain k8s-configmap target")
	assert.Contains(t, makefileStr, "k8s-deploy",
		"Makefile should contain k8s-deploy target")

	t.Logf("✓ Verified Makefile created with Kubernetes targets")
}
