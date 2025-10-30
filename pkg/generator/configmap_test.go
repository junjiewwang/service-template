package generator

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
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

			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				for _, check := range tt.checks {
					if !strings.Contains(content, check) {
						t.Errorf("Generated content missing expected string: %s", check)
					}
				}
			}
		})
	}
}
