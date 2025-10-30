package config

import (
	"testing"
)

func TestLoader_Load(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid configuration",
			yaml: `
service:
  name: test-service
  description: "Test Service"
  ports:
    - name: http
      port: 8080
      protocol: TCP
      expose: true
  deploy_dir: /usr/local/services

language:
  type: go
  version: "1.23"

build:
  dependency_files:
    auto_detect: true
  builder_image:
    amd64: "builder:amd64"
    arm64: "builder:arm64"
  runtime_image:
    amd64: "runtime:amd64"
    arm64: "runtime:arm64"
  commands:
    build: "go build"
  output_dir: dist

runtime:
  healthcheck:
    enabled: true
    type: http
    http:
      path: /health
      port: 8080
      timeout: 3
  startup:
    command: "./app"

local_dev:
  compose:
    volumes: []
  kubernetes:
    enabled: false
    namespace: default
    output_dir: k8s

metadata:
  template_version: "2.0.0"
  generator: "tcs-gen"
`,
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			yaml:    "invalid: yaml: content:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFromBytes([]byte(tt.yaml))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cfg == nil {
				t.Error("LoadFromBytes() returned nil config")
			}
			if !tt.wantErr && cfg.Service.Name != "test-service" {
				t.Errorf("LoadFromBytes() service name = %v, want test-service", cfg.Service.Name)
			}
		})
	}
}

func TestServiceConfig_Ports(t *testing.T) {
	yaml := `
service:
  name: test-service
  description: "Test"
  ports:
    - name: http
      port: 8080
      protocol: TCP
      expose: true
    - name: metrics
      port: 9090
      protocol: TCP
      expose: false
  deploy_dir: /usr/local/services

language:
  type: go
  version: "1.23"

build:
  dependency_files:
    auto_detect: true
  builder_image:
    amd64: "builder:amd64"
    arm64: "builder:arm64"
  runtime_image:
    amd64: "runtime:amd64"
    arm64: "runtime:arm64"
  commands:
    build: "go build"
  output_dir: dist

runtime:
  healthcheck:
    enabled: true
    type: http
    http:
      path: /health
      port: 8080
      timeout: 3
  startup:
    command: "./app"

local_dev:
  compose:
    volumes: []
  kubernetes:
    enabled: false
    namespace: default
    output_dir: k8s

metadata:
  template_version: "2.0.0"
  generator: "tcs-gen"
`

	cfg, err := LoadFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("LoadFromBytes() error = %v", err)
	}

	if len(cfg.Service.Ports) != 2 {
		t.Errorf("Expected 2 ports, got %d", len(cfg.Service.Ports))
	}

	if cfg.Service.Ports[0].Port != 8080 {
		t.Errorf("Expected first port to be 8080, got %d", cfg.Service.Ports[0].Port)
	}

	if cfg.Service.Ports[1].Port != 9090 {
		t.Errorf("Expected second port to be 9090, got %d", cfg.Service.Ports[1].Port)
	}
}
