package testutil

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
)

// TestNewBuilder 测试向后兼容的 NewBuilder API
func TestNewBuilder(t *testing.T) {
	cfg := NewBuilder().Build()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "test-service" {
		t.Errorf("expected service name 'test-service', got '%s'", cfg.Service.Name)
	}

	// 验证默认镜像配置
	if cfg.Build.BuilderImage != "@builders.go_1.21" {
		t.Errorf("expected builder image '@builders.go_1.21', got '%s'", cfg.Build.BuilderImage)
	}

	if cfg.Build.RuntimeImage != "@runtimes.alpine_3.18" {
		t.Errorf("expected runtime image '@runtimes.alpine_3.18', got '%s'", cfg.Build.RuntimeImage)
	}
}

// TestNewMinimal 测试向后兼容的 NewMinimal API
func TestNewMinimal(t *testing.T) {
	cfg := NewMinimal("my-service")

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "my-service" {
		t.Errorf("expected service name 'my-service', got '%s'", cfg.Service.Name)
	}
}

// TestNewConfigBuilder 测试新的 ConfigBuilder API
func TestNewConfigBuilder(t *testing.T) {
	cfg := NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuildCommand("go build -o bin/app").
		Build()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "test-service" {
		t.Errorf("expected service name 'test-service', got '%s'", cfg.Service.Name)
	}

	if cfg.Language.Type != "go" {
		t.Errorf("expected language 'go', got '%s'", cfg.Language.Type)
	}
}

// TestPresets 测试预设配置
func TestPresets(t *testing.T) {
	tests := []struct {
		name   string
		preset func() *config.ServiceConfig
	}{
		{"MinimalConfig", MinimalConfig},
		{"GoServiceConfig", GoServiceConfig},
		{"PythonServiceConfig", PythonServiceConfig},
		{"JavaServiceConfig", JavaServiceConfig},
		{"ConfigWithPlugins", ConfigWithPlugins},
		{"ConfigWithCustomHealthcheck", ConfigWithCustomHealthcheck},
		{"ConfigWithMultiArchPlugin", ConfigWithMultiArchPlugin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.preset()
			if cfg == nil {
				t.Fatalf("%s returned nil config", tt.name)
			}
			if cfg.Service.Name == "" {
				t.Errorf("%s: service name should not be empty", tt.name)
			}
		})
	}
}

// TestGeneratorTestConfig 测试生成器专用配置
func TestGeneratorTestConfig(t *testing.T) {
	cfg := NewGeneratorTestConfig()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "generator-test" {
		t.Errorf("expected service name 'generator-test', got '%s'", cfg.Service.Name)
	}

	if cfg.CI.ScriptDir != ".tad/build/generator-test" {
		t.Errorf("expected CI script dir '.tad/build/generator-test', got '%s'", cfg.CI.ScriptDir)
	}
}

// TestGeneratorTestConfigWithPlugins 测试带插件的生成器配置
func TestGeneratorTestConfigWithPlugins(t *testing.T) {
	cfg := NewGeneratorTestConfigWithPlugins()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if len(cfg.Plugins.Items) == 0 {
		t.Error("expected at least one plugin")
	}

	if cfg.Plugins.InstallDir != "/tce" {
		t.Errorf("expected plugin install dir '/tce', got '%s'", cfg.Plugins.InstallDir)
	}
}

// TestNewTestConfig 测试 fixtures 中的配置
func TestNewTestConfig(t *testing.T) {
	cfg := NewTestConfig()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "test-service" {
		t.Errorf("expected service name 'test-service', got '%s'", cfg.Service.Name)
	}
}

// TestNewTestConfigWithPlugins 测试 fixtures 中带插件的配置
func TestNewTestConfigWithPlugins(t *testing.T) {
	cfg := NewTestConfigWithPlugins()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if len(cfg.Plugins.Items) == 0 {
		t.Error("expected at least one plugin")
	}
}

// TestDefaultBaseImages 测试默认镜像配置
func TestDefaultBaseImages(t *testing.T) {
	images := DefaultBaseImages()

	if len(images.Builders) == 0 {
		t.Error("expected at least one builder image")
	}

	if len(images.Runtimes) == 0 {
		t.Error("expected at least one runtime image")
	}
}

// TestOptionsPattern 测试选项模式
func TestOptionsPattern(t *testing.T) {
	cfg := NewConfigWithOptions(
		MinimalConfig(),
		WithServiceNameOpt("option-service"),
		WithLanguageOpt("python"),
		WithBuildCommandOpt("pip install -r requirements.txt"),
	)

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "option-service" {
		t.Errorf("expected service name 'option-service', got '%s'", cfg.Service.Name)
	}

	if cfg.Language.Type != "python" {
		t.Errorf("expected language 'python', got '%s'", cfg.Language.Type)
	}
}
