package config

import (
	"strings"
	"testing"
)

// createTestBaseImages 创建测试用的基础镜像配置
func createTestBaseImages() BaseImagesConfig {
	return BaseImagesConfig{
		Builders: map[string]ArchImageConfig{
			"test_builder": {
				AMD64: "builder:amd64",
				ARM64: "builder:arm64",
			},
		},
		Runtimes: map[string]ArchImageConfig{
			"test_runtime": {
				AMD64: "runtime:amd64",
				ARM64: "runtime:arm64",
			},
		},
	}
}

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
				BaseImages: createTestBaseImages(),
				Service: ServiceInfo{
					Name:        "test-service",
					Description: "Test",
					Ports: []PortConfig{
						{Name: "http", Port: 8080, Protocol: "TCP", Expose: true},
					},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type: "go",
				},
				Build: BuildConfig{
					BuilderImage: "@builders.test_builder",
					RuntimeImage: "@runtimes.test_runtime",
					Commands: BuildCommandsConfig{
						Build: "go build",
					},
				},
				Runtime: RuntimeConfig{
					Healthcheck: HealthcheckConfig{
						Enabled: true,
						Type:    "default",
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
				BaseImages: createTestBaseImages(),
				Service: ServiceInfo{
					Name: "",
					Ports: []PortConfig{
						{Name: "http", Port: 8080, Protocol: "TCP"},
					},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type: "go",
				},
				Build: BuildConfig{
					BuilderImage: "@builders.test_builder",
					RuntimeImage: "@runtimes.test_runtime",
					Commands:     BuildCommandsConfig{Build: "build"},
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
				BaseImages: createTestBaseImages(),
				Service: ServiceInfo{
					Name: "test",
					Ports: []PortConfig{
						{Name: "http", Port: 99999, Protocol: "TCP"},
					},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type: "go",
				},
				Build: BuildConfig{
					BuilderImage: "@builders.test_builder",
					RuntimeImage: "@runtimes.test_runtime",
					Commands:     BuildCommandsConfig{Build: "build"},
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
				BaseImages: createTestBaseImages(),
				Service: ServiceInfo{
					Name:      "test",
					Ports:     []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type: "invalid",
				},
				Build: BuildConfig{
					BuilderImage: "@builders.test_builder",
					RuntimeImage: "@runtimes.test_runtime",
					Commands:     BuildCommandsConfig{Build: "build"},
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
				BaseImages: createTestBaseImages(),
				Service: ServiceInfo{
					Name:      "test",
					Ports:     []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{
					Type: "go",
				},
				Build: BuildConfig{
					BuilderImage: "@builders.test_builder",
					RuntimeImage: "@runtimes.test_runtime",
					Commands:     BuildCommandsConfig{Build: ""},
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
			name: "valid default healthcheck",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "default",
			},
			wantErr: false,
		},
		{
			name: "valid empty type defaults to default",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "",
			},
			wantErr: false,
		},
		{
			name: "valid custom healthcheck",
			healthcheck: HealthcheckConfig{
				Enabled:      true,
				Type:         "custom",
				CustomScript: "#!/bin/sh\nexit 0",
			},
			wantErr: false,
		},
		{
			name: "invalid healthcheck type",
			healthcheck: HealthcheckConfig{
				Enabled: true,
				Type:    "http",
			},
			wantErr: true,
			errMsg:  "is not valid",
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
				BaseImages: createTestBaseImages(),
				Service: ServiceInfo{
					Name:      "test",
					Ports:     []PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
					DeployDir: "/usr/local/services",
				},
				Language: LanguageConfig{Type: "go"},
				Build: BuildConfig{
					BuilderImage: "@builders.test_builder",
					RuntimeImage: "@runtimes.test_runtime",
					Commands:     BuildCommandsConfig{Build: "build"},
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
