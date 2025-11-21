package services

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependencyService_GetBuildDependencies(t *testing.T) {
	// Arrange
	cfg := testutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithDeployDir("/app").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithSystemPackages([]string{"git", "make", "gcc"}).
		WithCustomPackage(config.CustomPackage{
			Name:        "nacos",
			Description: "Nacos Service Discovery",
			InstallCommand: `echo "Installing to ${BUILD_OUTPUT_DIR}"
curl -L https://example.com/nacos.tar.gz -o ${BUILD_OUTPUT_DIR}/nacos.tar.gz`,
			Required: true,
		}).
		WithCustomPackage(config.CustomPackage{
			Name:           "consul",
			Description:    "Consul Service Mesh",
			InstallCommand: `curl -L https://example.com/consul.zip -o /tmp/consul.zip`,
			Required:       false,
		}).
		Build()

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	engine := core.NewTemplateEngine()
	svc := NewDependencyService(ctx, engine)

	// Act
	deps := svc.GetBuildDependencies()

	// Assert
	assert.Equal(t, 3, len(deps.SystemPkgs), "Should have 3 system packages")
	assert.Contains(t, deps.SystemPkgs, "git")
	assert.Contains(t, deps.SystemPkgs, "make")
	assert.Contains(t, deps.SystemPkgs, "gcc")

	assert.Equal(t, 2, len(deps.CustomPkgs), "Should have 2 custom packages")

	// Check first custom package
	assert.Equal(t, "nacos", deps.CustomPkgs[0].Name)
	assert.Equal(t, "Nacos Service Discovery", deps.CustomPkgs[0].Description)
	assert.True(t, deps.CustomPkgs[0].Required)
	// Variable should be replaced
	assert.Contains(t, deps.CustomPkgs[0].InstallCommand, "/opt/dist")

	// Check second custom package
	assert.Equal(t, "consul", deps.CustomPkgs[1].Name)
	assert.False(t, deps.CustomPkgs[1].Required)
}

func TestDependencyService_GetRuntimeDependencies(t *testing.T) {
	// Arrange
	cfg := testutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithRuntimePackages([]string{"ca-certificates", "tzdata"}).
		Build()

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	engine := core.NewTemplateEngine()
	svc := NewDependencyService(ctx, engine)

	// Act
	deps := svc.GetRuntimeDependencies()

	// Assert
	assert.Equal(t, 2, len(deps))
	assert.Contains(t, deps, "ca-certificates")
	assert.Contains(t, deps, "tzdata")
}

func TestDependencyService_HasBuildDependencies(t *testing.T) {
	tests := []struct {
		name       string
		systemPkgs []string
		customPkgs []config.CustomPackage
		want       bool
	}{
		{
			name:       "has system packages",
			systemPkgs: []string{"git"},
			customPkgs: nil,
			want:       true,
		},
		{
			name:       "has custom packages",
			systemPkgs: nil,
			customPkgs: []config.CustomPackage{
				{Name: "nacos", InstallCommand: "echo test", Required: true},
			},
			want: true,
		},
		{
			name:       "has both",
			systemPkgs: []string{"git"},
			customPkgs: []config.CustomPackage{
				{Name: "nacos", InstallCommand: "echo test", Required: true},
			},
			want: true,
		},
		{
			name:       "has none",
			systemPkgs: nil,
			customPkgs: nil,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := testutil.NewConfigBuilder().WithService("test", "Test Service")
			if len(tt.systemPkgs) > 0 {
				builder = builder.WithSystemPackages(tt.systemPkgs)
			}
			for _, pkg := range tt.customPkgs {
				builder = builder.WithCustomPackage(pkg)
			}
			cfg := builder.Build()

			ctx := context.NewGeneratorContext(cfg, "/tmp/output")
			engine := core.NewTemplateEngine()
			svc := NewDependencyService(ctx, engine)

			assert.Equal(t, tt.want, svc.HasBuildDependencies())
		})
	}
}

func TestDependencyService_HasSystemPackages(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("test", "Test Service").
		WithSystemPackages([]string{"git", "make"}).
		Build()

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	engine := core.NewTemplateEngine()
	svc := NewDependencyService(ctx, engine)

	assert.True(t, svc.HasSystemPackages())
}

func TestDependencyService_HasCustomPackages(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("test", "Test Service").
		WithCustomPackage(config.CustomPackage{
			Name:           "nacos",
			InstallCommand: "echo test",
			Required:       true,
		}).
		Build()

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	engine := core.NewTemplateEngine()
	svc := NewDependencyService(ctx, engine)

	assert.True(t, svc.HasCustomPackages())
}

func TestDependencyService_VariableSubstitution(t *testing.T) {
	// Arrange
	cfg := testutil.NewConfigBuilder().
		WithService("my-service", "My Service").
		WithDeployDir("/usr/local/services").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithCustomPackage(config.CustomPackage{
			Name: "test-pkg",
			InstallCommand: `echo "Service: ${SERVICE_NAME}"
echo "Deploy: ${DEPLOY_DIR}"
echo "Output: ${BUILD_OUTPUT_DIR}"
echo "Root: ${SERVICE_ROOT}"`,
			Required: true,
		}).
		Build()

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	engine := core.NewTemplateEngine()
	svc := NewDependencyService(ctx, engine)

	// Act
	deps := svc.GetBuildDependencies()

	// Assert
	require.Equal(t, 1, len(deps.CustomPkgs))
	installCmd := deps.CustomPkgs[0].InstallCommand

	// Debug output
	t.Logf("Install command after substitution:\n%s", installCmd)

	// Check that variables are replaced
	assert.Contains(t, installCmd, "my-service")
	assert.Contains(t, installCmd, "/usr/local/services")
	assert.Contains(t, installCmd, "/opt/dist")
	assert.Contains(t, installCmd, "/usr/local/services/my-service")

	// Check that variable placeholders are gone
	assert.NotContains(t, installCmd, "${SERVICE_NAME}")
	assert.NotContains(t, installCmd, "${DEPLOY_DIR}")
	assert.NotContains(t, installCmd, "${BUILD_OUTPUT_DIR}")
	assert.NotContains(t, installCmd, "${SERVICE_ROOT}")
}
