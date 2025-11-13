package compose

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Service.Ports = []config.PortConfig{
		{Port: 8080, Protocol: "tcp"},
		{Port: 9090, Protocol: "tcp"},
	}
	cfg.Runtime.Startup.Env = []config.EnvConfig{
		{Name: "ENV", Value: "production"},
	}
	cfg.LocalDev.Compose.Volumes = []config.VolumeConfig{
		{Source: "./data", Target: "${SERVICE_ROOT}/data"},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify content
	if !strings.Contains(content, "services:") {
		t.Error("Expected services field not found")
	}
	if !strings.Contains(content, "8080") {
		t.Error("Expected port 8080 not found")
	}
	if !strings.Contains(content, "9090") {
		t.Error("Expected port 9090 not found")
	}
	if !strings.Contains(content, "ENV=production") {
		t.Error("Expected environment variable not found")
	}
}

func TestGenerator_GetName(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if gen.GetName() != GeneratorType {
		t.Errorf("Expected name %s, got %s", GeneratorType, gen.GetName())
	}
}

func TestGenerator_Validate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if err := gen.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}

func TestGenerator_MergeEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name        string
		runtimeEnv  []config.EnvConfig
		composeEnv  []config.EnvConfig
		wantCount   int
		wantContain map[string]string
	}{
		{
			name: "Only runtime env",
			runtimeEnv: []config.EnvConfig{
				{Name: "ENV", Value: "production"},
				{Name: "LOG_LEVEL", Value: "info"},
			},
			composeEnv: []config.EnvConfig{},
			wantCount:  2,
			wantContain: map[string]string{
				"ENV":       "production",
				"LOG_LEVEL": "info",
			},
		},
		{
			name:       "Only compose env",
			runtimeEnv: []config.EnvConfig{},
			composeEnv: []config.EnvConfig{
				{Name: "DEBUG", Value: "true"},
				{Name: "LOG_LEVEL", Value: "debug"},
			},
			wantCount: 2,
			wantContain: map[string]string{
				"DEBUG":     "true",
				"LOG_LEVEL": "debug",
			},
		},
		{
			name: "Merge with override",
			runtimeEnv: []config.EnvConfig{
				{Name: "ENV", Value: "production"},
				{Name: "LOG_LEVEL", Value: "info"},
			},
			composeEnv: []config.EnvConfig{
				{Name: "LOG_LEVEL", Value: "debug"}, // Override runtime value
				{Name: "DEBUG", Value: "true"},      // New variable
			},
			wantCount: 3,
			wantContain: map[string]string{
				"ENV":       "production",
				"LOG_LEVEL": "debug", // Should be overridden by compose
				"DEBUG":     "true",
			},
		},
		{
			name:        "Empty both",
			runtimeEnv:  []config.EnvConfig{},
			composeEnv:  []config.EnvConfig{},
			wantCount:   0,
			wantContain: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testutil.NewTestConfig()
			cfg.Runtime.Startup.Env = tt.runtimeEnv
			cfg.LocalDev.Compose.Environment = tt.composeEnv

			ctx := context.NewGeneratorContext(cfg, "/tmp/output")
			gen, _ := New(ctx)
			generator := gen.(*Generator)

			result := generator.mergeEnvironmentVariables(ctx)

			// Check count
			if len(result) != tt.wantCount {
				t.Errorf("Expected %d env vars, got %d", tt.wantCount, len(result))
			}

			// Check content
			resultMap := make(map[string]string)
			for _, item := range result {
				envMap := item.(map[string]interface{})
				resultMap[envMap["Name"].(string)] = envMap["Value"].(string)
			}

			for name, expectedValue := range tt.wantContain {
				if actualValue, ok := resultMap[name]; !ok {
					t.Errorf("Expected env var %s not found", name)
				} else if actualValue != expectedValue {
					t.Errorf("Env var %s: expected value %s, got %s", name, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestGenerator_Generate_WithComposeEnvironment(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Service.Ports = []config.PortConfig{
		{Port: 8080, Protocol: "tcp"},
	}
	cfg.Runtime.Startup.Env = []config.EnvConfig{
		{Name: "ENV", Value: "production"},
		{Name: "LOG_LEVEL", Value: "info"},
	}
	cfg.LocalDev.Compose.Environment = []config.EnvConfig{
		{Name: "LOG_LEVEL", Value: "debug"}, // Override runtime value
		{Name: "DEBUG", Value: "true"},      // New variable
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify that compose environment variables are present
	if !strings.Contains(content, "DEBUG=true") {
		t.Error("Expected DEBUG=true not found in generated content")
	}

	// Verify that compose env overrides runtime env
	if !strings.Contains(content, "LOG_LEVEL=debug") {
		t.Error("Expected LOG_LEVEL=debug (overridden value) not found")
	}

	// Verify that runtime env is still present
	if !strings.Contains(content, "ENV=production") {
		t.Error("Expected ENV=production not found")
	}
}

func TestGenerator_Generate_WithEntrypoint(t *testing.T) {
	tests := []struct {
		name       string
		entrypoint []string
		wantLines  []string
	}{
		{
			name:       "Simple shell entrypoint",
			entrypoint: []string{"/bin/sh"},
			wantLines: []string{
				"entrypoint:",
				"- /bin/sh",
			},
		},
		{
			name: "Multi-line entrypoint with script",
			entrypoint: []string{
				"/bin/sh",
				"-c",
				"echo 'Starting service...' && exec /usr/local/services/test-service/bin/test-service",
			},
			wantLines: []string{
				"entrypoint:",
				"- /bin/sh",
				"- -c",
				"- echo 'Starting service...' && exec /usr/local/services/test-service/bin/test-service",
			},
		},
		{
			name:       "Empty entrypoint",
			entrypoint: []string{},
			wantLines:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testutil.NewTestConfig()
			cfg.Service.Ports = []config.PortConfig{
				{Port: 8080, Protocol: "tcp"},
			}
			cfg.LocalDev.Compose.Entrypoint = tt.entrypoint

			ctx := context.NewGeneratorContext(cfg, "/tmp/output")
			gen, err := New(ctx)
			if err != nil {
				t.Fatalf("Failed to create generator: %v", err)
			}

			content, err := gen.Generate()
			if err != nil {
				t.Fatalf("Failed to generate: %v", err)
			}

			// Verify expected lines
			for _, line := range tt.wantLines {
				if !strings.Contains(content, line) {
					t.Errorf("Expected line '%s' not found in generated content", line)
				}
			}

			// If entrypoint is empty, verify that entrypoint section is not present
			if len(tt.entrypoint) == 0 {
				if strings.Contains(content, "entrypoint:") {
					t.Error("Expected no entrypoint section, but found one")
				}
			}
		})
	}
}

func TestGenerator_Generate_WithEntrypointAndEnvironment(t *testing.T) {
	cfg := testutil.NewTestConfig()
	cfg.Service.Ports = []config.PortConfig{
		{Port: 8080, Protocol: "tcp"},
	}
	cfg.Runtime.Startup.Env = []config.EnvConfig{
		{Name: "ENV", Value: "production"},
	}
	cfg.LocalDev.Compose.Environment = []config.EnvConfig{
		{Name: "DEBUG", Value: "true"},
	}
	cfg.LocalDev.Compose.Entrypoint = []string{
		"/bin/sh",
		"-c",
		"echo 'Debug mode enabled' && exec /app/start.sh",
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify both environment and entrypoint are present
	if !strings.Contains(content, "environment:") {
		t.Error("Expected environment section not found")
	}
	if !strings.Contains(content, "DEBUG=true") {
		t.Error("Expected DEBUG=true not found")
	}
	if !strings.Contains(content, "entrypoint:") {
		t.Error("Expected entrypoint section not found")
	}
	if !strings.Contains(content, "/bin/sh") {
		t.Error("Expected /bin/sh in entrypoint not found")
	}
	if !strings.Contains(content, "echo 'Debug mode enabled' && exec /app/start.sh") {
		t.Error("Expected entrypoint command not found")
	}
}
