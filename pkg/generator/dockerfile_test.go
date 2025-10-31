package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerfileGenerator_Generate(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
			Ports: []config.PortConfig{
				{Name: "http", Port: 8080, Protocol: "TCP", Expose: true},
			},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{
			Type:    "go",
			Version: "1.23",
		},
		Build: config.BuildConfig{
			DependencyFiles: config.DependencyFilesConfig{
				AutoDetect: true,
			},
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.23-alpine",
				ARM64: "golang:1.23-alpine",
			},
			RuntimeImage: config.ArchImageConfig{
				AMD64: "alpine:latest",
				ARM64: "alpine:latest",
			},
			SystemDependencies: config.SystemDependenciesConfig{
				Build: config.PackagesConfig{
					Packages: []string{"git", "make"},
				},
			},
			Commands: config.BuildCommandsConfig{
				PreBuild:  "echo 'Pre-build'",
				Build:     "go build -o app",
				PostBuild: "echo 'Post-build'",
			},
			OutputDir: "dist",
		},
		Runtime: config.RuntimeConfig{
			SystemDependencies: config.SystemDependenciesConfig{
				Runtime: config.PackagesConfig{
					Packages: []string{"ca-certificates"},
				},
			},
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "http",
			},
		},
		Metadata: config.MetadataConfig{
			GeneratedAt: "2024-01-01T00:00:00Z",
		},
	}

	engine := NewTemplateEngine()
	vars := NewVariables(cfg)
	generator := NewDockerfileGenerator(cfg, engine, vars)

	tests := []struct {
		arch string
	}{
		{"amd64"},
		{"arm64"},
	}

	for _, tt := range tests {
		t.Run(tt.arch, func(t *testing.T) {
			content, err := generator.Generate(tt.arch)
			require.NoError(t, err, "Generate() should not return an error")

			// Check that content contains expected sections
			expectedSections := []string{
				"FROM",
				"WORKDIR",
				"COPY",
				"RUN",
				"ENTRYPOINT",
			}

			for _, section := range expectedSections {
				assert.Contains(t, content, section, "Generated Dockerfile should contain section: %s", section)
			}

			// Check architecture-specific image variables
			if tt.arch == "amd64" {
				assert.Contains(t, content, "BUILDER_IMAGE_X86", "Dockerfile should contain BUILDER_IMAGE_X86 variable")
				assert.Contains(t, content, "TLINUX_BASE_IMAGE_X86", "Dockerfile should contain TLINUX_BASE_IMAGE_X86 variable")
			} else if tt.arch == "arm64" {
				assert.Contains(t, content, "BUILDER_IMAGE_ARM", "Dockerfile should contain BUILDER_IMAGE_ARM variable")
				assert.Contains(t, content, "TLINUX_BASE_IMAGE_ARM", "Dockerfile should contain TLINUX_BASE_IMAGE_ARM variable")
			}
		})
	}
}

func TestGetDefaultDependencyFiles(t *testing.T) {
	tests := []struct {
		language string
		want     []string
	}{
		{"go", []string{"go.mod", "go.sum"}},
		{"python", []string{"requirements.txt"}},
		{"nodejs", []string{"package.json", "package-lock.json"}},
		{"java", []string{"pom.xml"}},
		{"unknown", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			got := getDefaultDependencyFiles(tt.language)
			assert.Equal(t, tt.want, got, "getDefaultDependencyFiles() should return expected files")
		})
	}
}

func TestDetectPackageManager(t *testing.T) {
	tests := []struct {
		image string
		want  string
	}{
		{"alpine:latest", "apk"},
		{"debian:bullseye", "apt-get"},
		{"ubuntu:22.04", "apt-get"},
		{"centos:7", "yum"},
		{"tencentos:3", "yum"},
		{"fedora:38", "dnf"},
		{"unknown:latest", "yum"},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			got := detectPackageManager(tt.image)
			assert.Equal(t, tt.want, got, "detectPackageManager(%s) should return expected package manager", tt.image)
		})
	}
}
