package generator

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
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
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			// Check that content contains expected sections
			expectedSections := []string{
				"FROM",
				"WORKDIR",
				"COPY",
				"RUN",
				"EXPOSE",
				"HEALTHCHECK",
				"CMD",
			}

			for _, section := range expectedSections {
				if !strings.Contains(content, section) {
					t.Errorf("Generated Dockerfile missing section: %s", section)
				}
			}

			// Check architecture-specific image variables
			if tt.arch == "amd64" {
				if !strings.Contains(content, "BUILDER_IMAGE_X86") {
					t.Error("Dockerfile should contain BUILDER_IMAGE_X86 variable")
				}
				if !strings.Contains(content, "TLINUX_BASE_IMAGE_X86") {
					t.Error("Dockerfile should contain TLINUX_BASE_IMAGE_X86 variable")
				}
			} else if tt.arch == "arm64" {
				if !strings.Contains(content, "BUILDER_IMAGE_ARM") {
					t.Error("Dockerfile should contain BUILDER_IMAGE_ARM variable")
				}
				if !strings.Contains(content, "TLINUX_BASE_IMAGE_ARM") {
					t.Error("Dockerfile should contain TLINUX_BASE_IMAGE_ARM variable")
				}
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
			if len(got) != len(tt.want) {
				t.Errorf("getDefaultDependencyFiles() = %v, want %v", got, tt.want)
				return
			}
			for i, file := range got {
				if file != tt.want[i] {
					t.Errorf("getDefaultDependencyFiles()[%d] = %v, want %v", i, file, tt.want[i])
				}
			}
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
			if got != tt.want {
				t.Errorf("detectPackageManager(%s) = %v, want %v", tt.image, got, tt.want)
			}
		})
	}
}
