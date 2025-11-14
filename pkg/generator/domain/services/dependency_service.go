package services

import (
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

// DependencyService handles dependency-related business logic
type DependencyService struct {
	ctx    *context.GeneratorContext
	engine *core.TemplateEngine
}

// NewDependencyService creates a new dependency service
func NewDependencyService(ctx *context.GeneratorContext, engine *core.TemplateEngine) *DependencyService {
	return &DependencyService{
		ctx:    ctx,
		engine: engine,
	}
}

// BuildDependency represents build dependencies (domain model)
type BuildDependency struct {
	SystemPkgs []string
	CustomPkgs []CustomPackageInfo
}

// CustomPackageInfo represents custom package information (domain model)
type CustomPackageInfo struct {
	Name           string
	Description    string
	InstallCommand string
	Required       bool
}

// GetBuildDependencies returns build dependencies with variable substitution
func (s *DependencyService) GetBuildDependencies() BuildDependency {
	deps := s.ctx.Config.Build.Dependencies

	// Get variables using the new variable system
	composer := s.ctx.GetVariableComposer().WithCommon().WithBuild()
	vars := s.convertToStringMap(composer.Build())

	// Process custom packages with variable substitution
	customPkgs := make([]CustomPackageInfo, len(deps.CustomPkgs))
	for i, pkg := range deps.CustomPkgs {
		// Use template engine for variable substitution
		installCmd := s.engine.ReplaceVariables(pkg.InstallCommand, vars)

		customPkgs[i] = CustomPackageInfo{
			Name:           pkg.Name,
			Description:    pkg.Description,
			InstallCommand: installCmd,
			Required:       pkg.Required,
		}
	}

	return BuildDependency{
		SystemPkgs: deps.SystemPkgs,
		CustomPkgs: customPkgs,
	}
}

// convertToStringMap converts map[string]interface{} to map[string]string
func (s *DependencyService) convertToStringMap(vars map[string]interface{}) map[string]string {
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}

// GetRuntimeDependencies returns runtime system dependencies
func (s *DependencyService) GetRuntimeDependencies() []string {
	return s.ctx.Config.Runtime.SystemDependencies.Packages
}

// HasBuildDependencies checks if there are any build dependencies
func (s *DependencyService) HasBuildDependencies() bool {
	deps := s.ctx.Config.Build.Dependencies
	return len(deps.SystemPkgs) > 0 || len(deps.CustomPkgs) > 0
}

// HasSystemPackages checks if there are system packages
func (s *DependencyService) HasSystemPackages() bool {
	return len(s.ctx.Config.Build.Dependencies.SystemPkgs) > 0
}

// HasCustomPackages checks if there are custom packages
func (s *DependencyService) HasCustomPackages() bool {
	return len(s.ctx.Config.Build.Dependencies.CustomPkgs) > 0
}

// GetSystemPackages returns system packages list
func (s *DependencyService) GetSystemPackages() []string {
	return s.ctx.Config.Build.Dependencies.SystemPkgs
}

// GetCustomPackages returns custom packages with variable substitution
func (s *DependencyService) GetCustomPackages() []CustomPackageInfo {
	return s.GetBuildDependencies().CustomPkgs
}
