package testutil_test

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/config/testutil"
)

// TestConfigBuilder 测试配置构建器
func TestConfigBuilder(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		Build()

	if cfg.Service.Name != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", cfg.Service.Name)
	}

	if cfg.Language.Type != "go" {
		t.Errorf("Expected language 'go', got '%s'", cfg.Language.Type)
	}

	if string(cfg.Build.BuilderImage) != "@builders.go_1.21" {
		t.Errorf("Expected builder image '@builders.go_1.21', got '%s'", cfg.Build.BuilderImage)
	}
}

// TestPresets 测试预设配置
func TestPresets(t *testing.T) {
	tests := []struct {
		name   string
		preset func() *config.ServiceConfig
	}{
		{"MinimalConfig", testutil.MinimalConfig},
		{"GoServiceConfig", testutil.GoServiceConfig},
		{"PythonServiceConfig", testutil.PythonServiceConfig},
		{"JavaServiceConfig", testutil.JavaServiceConfig},
		{"ConfigWithPlugins", testutil.ConfigWithPlugins},
		{"ConfigWithCustomHealthcheck", testutil.ConfigWithCustomHealthcheck},
		{"ConfigWithMultiArchPlugin", testutil.ConfigWithMultiArchPlugin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.preset()
			if cfg == nil {
				t.Fatal("Preset returned nil config")
			}

			// 验证配置
			validator := config.NewValidator(cfg)
			if err := validator.Validate(); err != nil {
				t.Errorf("Preset %s validation failed: %v", tt.name, err)
			}
		})
	}
}

// TestOptions 测试配置选项
func TestOptions(t *testing.T) {
	cfg := testutil.NewConfigWithOptions(
		testutil.MinimalConfig(),
		testutil.WithServiceNameOpt("custom-service"),
		testutil.WithPortOpt("http", 8080, "TCP", true),
		testutil.WithBuildCommandOpt("make build"),
	)

	if cfg.Service.Name != "custom-service" {
		t.Errorf("Expected service name 'custom-service', got '%s'", cfg.Service.Name)
	}

	if len(cfg.Service.Ports) == 0 {
		t.Fatal("Expected at least one port")
	}

	if cfg.Service.Ports[0].Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Service.Ports[0].Port)
	}

	if cfg.Build.Commands.Build != "make build" {
		t.Errorf("Expected build command 'make build', got '%s'", cfg.Build.Commands.Build)
	}
}

// TestCombinedPatterns 测试组合模式
func TestCombinedPatterns(t *testing.T) {
	// 从预设开始
	cfg := testutil.GoServiceConfig()

	// 使用选项修改
	cfg = testutil.ApplyOptions(cfg,
		testutil.WithServiceNameOpt("my-go-service"),
		testutil.WithPortOpt("grpc", 9000, "TCP", true),
	)

	if cfg.Service.Name != "my-go-service" {
		t.Errorf("Expected service name 'my-go-service', got '%s'", cfg.Service.Name)
	}

	// 验证原有端口仍然存在
	foundHTTP := false
	foundGRPC := false
	for _, port := range cfg.Service.Ports {
		if port.Name == "http" {
			foundHTTP = true
		}
		if port.Name == "grpc" {
			foundGRPC = true
		}
	}

	if !foundHTTP {
		t.Error("Expected to find http port")
	}
	if !foundGRPC {
		t.Error("Expected to find grpc port")
	}
}

// TestBuilderWithDefaults 测试带默认值的构建
func TestBuilderWithDefaults(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		BuildWithDefaults()

	// 验证默认值已填充
	if cfg.Service.DeployDir == "" {
		t.Error("Expected deploy_dir to be filled with default value")
	}

	if cfg.Metadata.TemplateVersion == "" {
		t.Error("Expected template_version to be filled with default value")
	}

	if cfg.Metadata.Generator == "" {
		t.Error("Expected generator to be filled with default value")
	}
}

// BenchmarkConfigBuilder 性能测试
func BenchmarkConfigBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = testutil.NewConfigBuilder().
			WithService("test-service", "Test Service").
			WithLanguage("go").
			WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
			WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
			WithBuilderImage("@builders.go_1.21").
			WithRuntimeImage("@runtimes.alpine_3.18").
			WithBuildCommand("go build -o bin/app").
			Build()
	}
}

// BenchmarkPreset 预设性能测试
func BenchmarkPreset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = testutil.GoServiceConfig()
	}
}
