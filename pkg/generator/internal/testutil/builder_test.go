package testutil

import (
	"testing"
)

func TestNewMinimalConfig(t *testing.T) {
	cfg := NewMinimalConfig().Build()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "test-service" {
		t.Errorf("expected service name 'test-service', got '%s'", cfg.Service.Name)
	}

	if cfg.Build.BuilderImage != "@builders.builder" {
		t.Errorf("expected builder image '@builders.builder', got '%s'", cfg.Build.BuilderImage)
	}

	if cfg.Build.RuntimeImage != "@runtimes.runtime" {
		t.Errorf("expected runtime image '@runtimes.runtime', got '%s'", cfg.Build.RuntimeImage)
	}
}

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig().Build()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Service.Name != "test-service" {
		t.Errorf("expected service name 'test-service', got '%s'", cfg.Service.Name)
	}

	if len(cfg.Service.Ports) == 0 {
		t.Error("expected at least one port")
	}

	if cfg.Build.BuilderImage != "@builders.test_builder" {
		t.Errorf("expected builder image '@builders.test_builder', got '%s'", cfg.Build.BuilderImage)
	}
}

func TestConfigBuilder_WithService(t *testing.T) {
	cfg := NewMinimalConfig().
		WithService(func(s *ServiceBuilder) {
			s.Name("custom-service").
				Description("Custom description").
				DeployDir("/custom/path").
				AddPort(9090, "grpc")
		}).
		Build()

	if cfg.Service.Name != "custom-service" {
		t.Errorf("expected service name 'custom-service', got '%s'", cfg.Service.Name)
	}

	if cfg.Service.Description != "Custom description" {
		t.Errorf("expected description 'Custom description', got '%s'", cfg.Service.Description)
	}

	if cfg.Service.DeployDir != "/custom/path" {
		t.Errorf("expected deploy dir '/custom/path', got '%s'", cfg.Service.DeployDir)
	}

	if len(cfg.Service.Ports) == 0 {
		t.Fatal("expected at least one port")
	}

	if cfg.Service.Ports[0].Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Service.Ports[0].Port)
	}
}

func TestConfigBuilder_WithBuild(t *testing.T) {
	cfg := NewMinimalConfig().
		WithBuild(func(b *BuildConfigBuilder) {
			b.BuildCommand("custom build").
				PreBuildCommand("custom pre-build").
				PostBuildCommand("custom post-build").
				AddSystemPackage("git").
				AddSystemPackage("curl")
		}).
		Build()

	if cfg.Build.Commands.Build != "custom build" {
		t.Errorf("expected build command 'custom build', got '%s'", cfg.Build.Commands.Build)
	}

	if cfg.Build.Commands.PreBuild != "custom pre-build" {
		t.Errorf("expected pre-build command 'custom pre-build', got '%s'", cfg.Build.Commands.PreBuild)
	}

	if cfg.Build.Commands.PostBuild != "custom post-build" {
		t.Errorf("expected post-build command 'custom post-build', got '%s'", cfg.Build.Commands.PostBuild)
	}

	if len(cfg.Build.Dependencies.SystemPkgs) != 2 {
		t.Errorf("expected 2 system packages, got %d", len(cfg.Build.Dependencies.SystemPkgs))
	}
}

func TestConfigBuilder_WithPlugins(t *testing.T) {
	cfg := NewMinimalConfig().
		WithPlugins(func(p *PluginBuilder) {
			p.InstallDir("/custom/tce").
				AddPlugin("plugin1", "Plugin 1", "https://example.com/plugin1").
				AddPlugin("plugin2", "Plugin 2", "https://example.com/plugin2")
		}).
		Build()

	if cfg.Plugins.InstallDir != "/custom/tce" {
		t.Errorf("expected install dir '/custom/tce', got '%s'", cfg.Plugins.InstallDir)
	}

	if len(cfg.Plugins.Items) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(cfg.Plugins.Items))
	}

	if cfg.Plugins.Items[0].Name != "plugin1" {
		t.Errorf("expected plugin name 'plugin1', got '%s'", cfg.Plugins.Items[0].Name)
	}

	if cfg.Plugins.Items[1].Name != "plugin2" {
		t.Errorf("expected plugin name 'plugin2', got '%s'", cfg.Plugins.Items[1].Name)
	}
}

func TestConfigBuilder_WithLanguage(t *testing.T) {
	cfg := NewMinimalConfig().
		WithLanguage("python").
		Build()

	if cfg.Language.Type != "python" {
		t.Errorf("expected language 'python', got '%s'", cfg.Language.Type)
	}
}

func TestConfigBuilder_WithBaseImagesPreset(t *testing.T) {
	cfg := NewMinimalConfig().
		WithBaseImagesPreset(CustomBaseImages).
		Build()

	if len(cfg.BaseImages.Builders) != 2 {
		t.Errorf("expected 2 builders, got %d", len(cfg.BaseImages.Builders))
	}

	if len(cfg.BaseImages.Runtimes) != 2 {
		t.Errorf("expected 2 runtimes, got %d", len(cfg.BaseImages.Runtimes))
	}
}

func TestNewGoServicePreset(t *testing.T) {
	cfg := NewGoServicePreset().Build()

	if cfg.Language.Type != "go" {
		t.Errorf("expected language 'go', got '%s'", cfg.Language.Type)
	}

	if cfg.Build.Commands.PreBuild != "go mod download" {
		t.Errorf("expected pre-build 'go mod download', got '%s'", cfg.Build.Commands.PreBuild)
	}
}

func TestNewPythonServicePreset(t *testing.T) {
	cfg := NewPythonServicePreset().Build()

	if cfg.Language.Type != "python" {
		t.Errorf("expected language 'python', got '%s'", cfg.Language.Type)
	}

	if cfg.Build.Commands.Build != "pip install -r requirements.txt" {
		t.Errorf("expected build command 'pip install -r requirements.txt', got '%s'", cfg.Build.Commands.Build)
	}
}

func TestNewServiceWithPluginsPreset(t *testing.T) {
	cfg := NewServiceWithPluginsPreset().Build()

	if cfg.Plugins.InstallDir != "/tce" {
		t.Errorf("expected install dir '/tce', got '%s'", cfg.Plugins.InstallDir)
	}
}

func TestNewConfigWithOptions(t *testing.T) {
	cfg := NewConfigWithOptions(
		WithServiceName("option-service"),
		WithDeployDir("/opt/option"),
		WithPort(8080, "http"),
		WithLanguageType("go"),
		WithBuildCommand("go build"),
	)

	if cfg.Service.Name != "option-service" {
		t.Errorf("expected service name 'option-service', got '%s'", cfg.Service.Name)
	}

	if cfg.Service.DeployDir != "/opt/option" {
		t.Errorf("expected deploy dir '/opt/option', got '%s'", cfg.Service.DeployDir)
	}

	if cfg.Language.Type != "go" {
		t.Errorf("expected language 'go', got '%s'", cfg.Language.Type)
	}

	if cfg.Build.Commands.Build != "go build" {
		t.Errorf("expected build command 'go build', got '%s'", cfg.Build.Commands.Build)
	}
}

func TestBaseImagesPreset_ToConfig(t *testing.T) {
	preset := MinimalBaseImages
	cfg := preset.ToConfig()

	if len(cfg.Builders) != 1 {
		t.Errorf("expected 1 builder, got %d", len(cfg.Builders))
	}

	if len(cfg.Runtimes) != 1 {
		t.Errorf("expected 1 runtime, got %d", len(cfg.Runtimes))
	}

	builder, ok := cfg.Builders["builder"]
	if !ok {
		t.Fatal("expected builder 'builder' to exist")
	}

	if builder.AMD64 != "builder:test" {
		t.Errorf("expected AMD64 image 'builder:test', got '%s'", builder.AMD64)
	}
}

func TestChainedBuilders(t *testing.T) {
	// 测试链式调用
	cfg := NewMinimalConfig().
		WithService(func(s *ServiceBuilder) {
			s.Name("chained-service").
				AddPort(8080, "http").
				AddPort(9090, "grpc")
		}).
		WithBuild(func(b *BuildConfigBuilder) {
			b.BuildCommand("make build").
				PreBuildCommand("make deps")
		}).
		WithPlugins(func(p *PluginBuilder) {
			p.AddPlugin("plugin1", "desc1", "url1")
		}).
		WithLanguage("go").
		Build()

	if cfg.Service.Name != "chained-service" {
		t.Errorf("expected service name 'chained-service', got '%s'", cfg.Service.Name)
	}

	if len(cfg.Service.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(cfg.Service.Ports))
	}

	if cfg.Build.Commands.Build != "make build" {
		t.Errorf("expected build command 'make build', got '%s'", cfg.Build.Commands.Build)
	}

	if len(cfg.Plugins.Items) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(cfg.Plugins.Items))
	}

	if cfg.Language.Type != "go" {
		t.Errorf("expected language 'go', got '%s'", cfg.Language.Type)
	}
}
