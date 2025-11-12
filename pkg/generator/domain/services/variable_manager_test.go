package services

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/domain/models"
)

func TestNewVariableManager(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "testservice",
			DeployDir: "/data/services",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/plugins",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	manager := NewVariableManager(ctx)

	if manager == nil {
		t.Fatal("NewVariableManager returned nil")
	}

	if manager.GetContext() != ctx {
		t.Error("GetContext() did not return the correct context")
	}

	pathsConfig := manager.GetPathsConfig()
	if pathsConfig.PluginInstallDir != "/plugins" {
		t.Errorf("PluginInstallDir = %v, want /plugins", pathsConfig.PluginInstallDir)
	}

	if pathsConfig.ServiceDeployDir != "/data/services" {
		t.Errorf("ServiceDeployDir = %v, want /data/services", pathsConfig.ServiceDeployDir)
	}
}

func TestNewVariableManagerWithPaths(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "testservice",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")

	customPaths := models.NewPathsConfig().
		WithPluginInstallDir("/custom/plugins").
		WithServiceDeployDir("/custom/services")

	manager := NewVariableManagerWithPaths(ctx, customPaths)

	pathsConfig := manager.GetPathsConfig()
	if pathsConfig.PluginInstallDir != "/custom/plugins" {
		t.Errorf("PluginInstallDir = %v, want /custom/plugins", pathsConfig.PluginInstallDir)
	}

	if pathsConfig.ServiceDeployDir != "/custom/services" {
		t.Errorf("ServiceDeployDir = %v, want /custom/services", pathsConfig.ServiceDeployDir)
	}
}

func TestVariableManager_PrepareForDockerfile(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "testservice",
		},
		Build: config.BuildConfig{
			BuilderImage: config.ArchImageConfig{
				AMD64: "golang:1.21-alpine",
			},
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	manager := NewVariableManager(ctx)

	composer := manager.PrepareForDockerfile("amd64")

	if composer == nil {
		t.Fatal("PrepareForDockerfile returned nil")
	}

	vars := composer.Build()
	if _, ok := vars["BUILDER_IMAGE"]; !ok {
		t.Error("PrepareForDockerfile did not include BUILDER_IMAGE")
	}
}

func TestVariableManager_PrepareForCompose(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "testservice",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	manager := NewVariableManager(ctx)

	composer := manager.PrepareForCompose()

	if composer == nil {
		t.Fatal("PrepareForCompose returned nil")
	}

	vars := composer.Build()
	if _, ok := vars["SERVICE_NAME"]; !ok {
		t.Error("PrepareForCompose did not include SERVICE_NAME")
	}
}

func TestVariableManager_PrepareForBuildScript(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "testservice",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	manager := NewVariableManager(ctx)

	composer := manager.PrepareForBuildScript()

	if composer == nil {
		t.Fatal("PrepareForBuildScript returned nil")
	}

	vars := composer.Build()
	if _, ok := vars["SERVICE_NAME"]; !ok {
		t.Error("PrepareForBuildScript did not include SERVICE_NAME")
	}
}

func TestVariableManager_PrepareForScript(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "testservice",
		},
		Language: config.LanguageConfig{
			Type: "go",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	manager := NewVariableManager(ctx)

	composer := manager.PrepareForScript()

	if composer == nil {
		t.Fatal("PrepareForScript returned nil")
	}

	vars := composer.Build()
	if _, ok := vars["LANGUAGE"]; !ok {
		t.Error("PrepareForScript did not include LANGUAGE")
	}
}

func TestVariableManager_AddPathVariables(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "testservice",
			DeployDir: "/data/services",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/plugins",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	manager := NewVariableManager(ctx)

	composer := ctx.GetVariablePreset().ForScript()
	composer = manager.AddPathVariables(composer, "testservice")

	vars := composer.Build()

	tests := []struct {
		key      string
		expected string
	}{
		{"SERVICE_ROOT", "/data/services/testservice"},
		{"SERVICE_BIN_PATH", "/data/services/testservice/bin"},
		{"SERVICE_CONFIG_PATH", "/data/services/testservice/conf"},
		{"SERVICE_LOG_PATH", "/data/services/testservice/logs"},
		{"SERVICE_DATA_PATH", "/data/services/testservice/data"},
		{"PLUGIN_INSTALL_DIR", "/plugins"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if val, ok := vars[tt.key]; !ok {
				t.Errorf("Variable %s not found", tt.key)
			} else if val != tt.expected {
				t.Errorf("%s = %v, want %v", tt.key, val, tt.expected)
			}
		})
	}
}

func TestVariableManager_PrepareWithPaths(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "testservice",
			DeployDir: "/data/services",
		},
		Plugins: config.PluginsConfig{
			InstallDir: "/plugins",
		},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	manager := NewVariableManager(ctx)

	composer := manager.PrepareWithPaths(func() *context.VariableComposer {
		return manager.PrepareForScript()
	})

	vars := composer.Build()

	// Check that both preset variables and path variables are present
	if _, ok := vars["SERVICE_NAME"]; !ok {
		t.Error("PrepareWithPaths did not include preset variables")
	}

	if _, ok := vars["SERVICE_ROOT"]; !ok {
		t.Error("PrepareWithPaths did not include path variables")
	}
}
