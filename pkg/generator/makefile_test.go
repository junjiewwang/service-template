package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakefileGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.ServiceConfig
		wantErr bool
		checks  []string
	}{
		{
			name: "basic makefile",
			config: &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/opt/services",
					Ports: []config.PortConfig{
						{Port: 8080, Protocol: "tcp"},
					},
				},
				Build: config.BuildConfig{},
				Language: config.LanguageConfig{
					Type:    "golang",
					Version: "1.21",
				},
				LocalDev: config.LocalDevConfig{
					Kubernetes: config.KubernetesConfig{
						Namespace:  "default",
						VolumeType: "configMap",
					},
				},
			},
			wantErr: false,
			checks: []string{
				"PROJECT_NAME ?= test-service",
				".PHONY:",
				"ARCH := $(shell uname -m)",
				"check-tools:",
				"k8s-convert:",
				"k8s-deploy:",
			},
		},
		{
			name: "makefile with service info",
			config: &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "multi-arch-service",
					DeployDir: "/opt/services",
					Ports: []config.PortConfig{
						{Port: 9090, Protocol: "tcp"},
					},
				},
				Build: config.BuildConfig{},
				Language: config.LanguageConfig{
					Type:    "golang",
					Version: "1.21",
				},
				LocalDev: config.LocalDevConfig{
					Kubernetes: config.KubernetesConfig{
						Namespace:  "production",
						VolumeType: "persistentVolumeClaim",
						OutputDir:  "k8s-output",
					},
				},
			},
			wantErr: false,
			checks: []string{
				"PROJECT_NAME ?= multi-arch-service",
				"K8S_NAMESPACE ?= production",
				"K8S_OUTPUT_DIR ?= k8s-output",
				"K8S_VOLUME_TYPE ?= persistentVolumeClaim",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup generator with configuration
			engine := NewTemplateEngine()
			vars := NewVariables(tt.config)
			g := NewMakefileGenerator(tt.config, engine, vars)
			require.NotNil(t, g, "Makefile generator should be created")

			// Act: Generate Makefile
			content, err := g.Generate()

			// Assert: Check results
			if tt.wantErr {
				assert.Error(t, err, "Generate() should return an error")
				t.Logf("Expected error occurred: %v", err)
			} else {
				require.NoError(t, err, "Generate() should not return an error")
				require.NotEmpty(t, content, "Generated Makefile should not be empty")

				t.Logf("Generated Makefile: %d bytes", len(content))

				// Verify all expected checks
				for i, check := range tt.checks {
					assert.Contains(t, content, check,
						"Generated Makefile should contain expected string [%d]: %s", i+1, check)
				}
				t.Logf("✓ Verified all %d expected sections present", len(tt.checks))
			}
		})
	}
}
