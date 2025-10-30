package config

import (
	"strings"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ServiceConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: &ServiceConfig{
				Service: ServiceInfo{
					Name:        "test-service",
					Description: "Test",
					Ports: []PortConfig{
						{Name: "http", Port: 8080, Protocol: "TCP", Expose: true},
					},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type:    "go",
					Version: "1.23",
				},
				Build: BuildConfig{
					BuilderImage: ArchImageConfig{
						AMD64: "builder:amd64",
						ARM64: "builder:arm64",
					},
					RuntimeImage: ArchImageConfig{
						AMD64: "runtime:amd64",
						ARM64: "runtime:arm64",
					},
					Commands: BuildCommandsConfig{
						Build: "go build",
					},
					OutputDir: "dist",
				},
				Runtime: RuntimeConfig{
					Healthcheck: HealthcheckConfig{
						Enabled: true,
						Type:    "http",
						HTTP: HTTPHealthConfig{
							Path:    "/health",
							Port:    8080,
							Timeout: 3,
						},
					},
					Startup: StartupConfig{
						Command: "./app",
					},
				},
				LocalDev: LocalDevConfig{
					Kubernetes: KubernetesConfig{
						Enabled:   false,
						Namespace: "default",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing service name",
			config: &ServiceConfig{
				Service: ServiceInfo{
					Name: "",
					Ports: []PortConfig{
						{Name: "http", Port: 8080, Protocol: "TCP"},
					},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type:    "go",
					Version: "1.23",
				},
				Build: BuildConfig{
					BuilderImage: ArchImageConfig{AMD64: "b", ARM64: "b"},
					RuntimeImage: ArchImageConfig{AMD64: "r", ARM64: "r"},
					Commands:     BuildCommandsConfig{Build: "build"},
					OutputDir:    "dist",
				},
				Runtime: RuntimeConfig{
					Startup: StartupConfig{Command: "./app"},
				},
			},
			wantErr: true,
			errMsg:  "service.name is required",
		},
		{
			name: "invalid port",
			config: &ServiceConfig{
				Service: ServiceInfo{
					Name: "test",
					Ports: []PortConfig{
						{Name: "http", Port: 99999, Protocol: "TCP"},
					},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type:    "go",
					Version: "1.23",
				},
				Build: BuildConfig{
					BuilderImage: ArchImageConfig{AMD64: "b", ARM64: "b"},
					RuntimeImage: ArchImageConfig{AMD64: "r", ARM64: "r"},
					Commands:     BuildCommandsConfig{Build: "build"},
					OutputDir:    "dist",
				},
				Runtime: RuntimeConfig{
					Startup: StartupConfig{Command: "./app"},
				},
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "invalid language type",
			config: &ServiceConfig{
				Service: ServiceInfo{
					Name:      "test",
					Ports:     []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type:    "invalid",
					Version: "1.0",
				},
				Build: BuildConfig{
					BuilderImage: ArchImageConfig{AMD64: "b", ARM64: "b"},
					RuntimeImage: ArchImageConfig{AMD64: "r", ARM64: "r"},
					Commands:     BuildCommandsConfig{Build: "build"},
					OutputDir:    "dist",
				},
				Runtime: RuntimeConfig{
					Startup: StartupConfig{Command: "./app"},
				},
			},
			wantErr: true,
			errMsg:  "is not supported",
		},
		{
			name: "missing build command",
			config: &ServiceConfig{
				Service: ServiceInfo{
					Name:      "test",
					Ports:     []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type:    "go",
					Version: "1.23",
				},
				Build: BuildConfig{
					BuilderImage: ArchImageConfig{AMD64: "b", ARM64: "b"},
					RuntimeImage: ArchImageConfig{AMD64: "r", ARM64: "r"},
					Commands:     BuildCommandsConfig{Build: ""},
					OutputDir:    "dist",
				},
				Runtime: RuntimeConfig{
					Startup: StartupConfig{Command: "./app"},
				},
			},
			wantErr: true,
			errMsg:  "build.commands.build is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.config)
			err := validator.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_ValidateHealthcheck(t *testing.T) {
	tests := []struct {
		name        string
		healthcheck HealthcheckConfig
		wantErr     bool
		errMsg      string
	}{
		{
			name: "valid http healthcheck",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "http",
				HTTP: HTTPHealthConfig{
					Path:    "/health",
					Port:    8080,
					Timeout: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid healthcheck type",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "invalid",
			},
			wantErr: true,
			errMsg:  "is not valid",
		},
		{
			name: "http healthcheck missing path",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "http",
				HTTP: HTTPHealthConfig{
					Port: 8080,
				},
			},
			wantErr: true,
			errMsg:  "path is required",
		},
		{
			name: "custom healthcheck missing script",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "custom",
			},
			wantErr: true,
			errMsg:  "custom_script is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ServiceConfig{
				Service: ServiceInfo{
					Name:      "test",
					Ports:     []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{Type: "go", Version: "1.23"},
				Build: BuildConfig{
					BuilderImage: ArchImageConfig{AMD64: "b", ARM64: "b"},
					RuntimeImage: ArchImageConfig{AMD64: "r", ARM64: "r"},
					Commands:     BuildCommandsConfig{Build: "build"},
					OutputDir:    "dist",
				},
				Runtime: RuntimeConfig{
					Healthcheck: tt.healthcheck,
					Startup:     StartupConfig{Command: "./app"},
				},
			}

			validator := NewValidator(config)
			err := validator.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}
