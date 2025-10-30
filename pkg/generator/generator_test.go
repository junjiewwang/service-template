package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
)

func TestGenerator_Generate(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
		LocalDev: config.LocalDevConfig{
			Kubernetes: config.KubernetesConfig{
				Enabled: true,
				ConfigMap: config.ConfigMapConfig{
					AutoDetect: true,
				},
			},
		},
		Plugins: []config.PluginConfig{
			{
				Name:        "test-plugin",
				Description: "Test plugin",
				DownloadURL: "https://example.com/plugin.tar.gz",
				InstallDir:  "/opt/plugins",
			},
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)
	err = gen.Generate()

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that expected files were created
	expectedFiles := []string{
		".tad/build/test-service/Dockerfile.test-service.amd64",
		".tad/build/test-service/Dockerfile.test-service.arm64",
		"compose.yaml",
		"Makefile",
		"configmap.yaml",
		"bk-ci/tcs/build.sh",
		"bk-ci/tcs/deps_install.sh",
		"bk-ci/tcs/rt_prepare.sh",
		"hooks/start.sh",
		"hooks/healthchk.sh",
		".tad/devops.yaml",
	}

	// List all generated files for debugging
	filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			relPath, _ := filepath.Rel(outputDir, path)
			t.Logf("Generated file: %s", relPath)
		}
		return nil
	})

	for _, file := range expectedFiles {
		fullPath := filepath.Join(outputDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s (full path: %s)", file, fullPath)
		}
	}
}

func TestGenerator_GenerateDockerfiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)

	// Call Generate which internally calls generateDockerfiles
	err = gen.Generate()

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check Dockerfile was created
	dockerfilePath := filepath.Join(outputDir, ".tad", "build", "test-service", "Dockerfile.test-service.amd64")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		t.Error("Dockerfile.test-service.amd64 not created")
	}
}

func TestGenerator_GenerateCompose(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)
	err = gen.Generate()

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check compose file was created
	composePath := filepath.Join(outputDir, "compose.yaml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		t.Error("compose.yaml not created")
	}
}

func TestGenerator_GenerateMakefile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)
	err = gen.Generate()

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check Makefile was created
	makefilePath := filepath.Join(outputDir, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		t.Error("Makefile not created")
	}
}

func TestGenerator_GenerateScripts(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)
	err = gen.Generate()

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check scripts were created
	expectedScripts := []string{
		"bk-ci/tcs/build.sh",
		"bk-ci/tcs/deps_install.sh",
		"bk-ci/tcs/rt_prepare.sh",
		"hooks/start.sh",
		"hooks/healthchk.sh",
	}

	for _, script := range expectedScripts {
		scriptPath := filepath.Join(outputDir, script)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			t.Errorf("Script not created: %s", script)
		}
	}
}

func TestGenerator_GenerateConfigMap(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/opt/services",
			Ports: []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			},
		},
		Build: config.BuildConfig{
			OutputDir: "build",
		},
		Language: config.LanguageConfig{
			Type:    "golang",
			Version: "1.21",
		},
		LocalDev: config.LocalDevConfig{
			Kubernetes: config.KubernetesConfig{
				Enabled: true,
				ConfigMap: config.ConfigMapConfig{
					AutoDetect: true,
				},
			},
		},
	}

	outputDir := filepath.Join(tmpDir, "output")
	gen := NewGenerator(cfg, outputDir)
	err = gen.Generate()

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check configmap was created (it's generated in root, not k8s-manifests)
	configmapPath := filepath.Join(outputDir, "configmap.yaml")
	if _, err := os.Stat(configmapPath); os.IsNotExist(err) {
		t.Error("configmap.yaml not created")
	}
}
