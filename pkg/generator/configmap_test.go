package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigMapGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.ServiceConfig
		wantErr bool
		checks  []string
	}{
		{
			name: "basic configmap",
			config: &config.ServiceConfig{
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
			},
			wantErr: false,
			checks: []string{
				"apiVersion: v1",
				"kind: ConfigMap",
				"name: test-service-config",
			},
		},
		{
			name: "configmap with custom config",
			config: &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "custom-service",
					DeployDir: "/opt/services",
					Ports: []config.PortConfig{
						{Port: 9090, Protocol: "tcp"},
					},
				},
				Build: config.BuildConfig{
					OutputDir: "build",
				},
				Language: config.LanguageConfig{
					Type:    "golang",
					Version: "1.21",
					Config: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			wantErr: false,
			checks: []string{
				"apiVersion: v1",
				"kind: ConfigMap",
				"name: custom-service-config",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewTemplateEngine()
			vars := NewVariables(tt.config)
			g := NewConfigMapGenerator(tt.config, engine, vars)
			content, err := g.Generate()

			if tt.wantErr {
				assert.Error(t, err, "Generate() should return an error")
			} else {
				require.NoError(t, err, "Generate() should not return an error")
				assert.NotEmpty(t, content, "Generated content should not be empty")

				for _, check := range tt.checks {
					assert.Contains(t, content, check, "Generated content should contain expected string: %s", check)
				}
			}
		})
	}
}
