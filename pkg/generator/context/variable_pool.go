package context

import (
	"sync"
)

// VariablePool manages shared template variables (Flyweight Pool)
// It creates and caches shared variable sets to avoid duplication
type VariablePool struct {
	ctx   *GeneratorContext
	cache map[string]*SharedVariables // Cache for different types of shared variables
	mu    sync.RWMutex
}

// SharedVariables represents a shareable set of variables (Flyweight Object)
type SharedVariables struct {
	category string                 // Variable category: common, build, runtime, plugin, etc.
	vars     map[string]interface{} // Actual variables
	frozen   bool                   // Whether frozen (immutable)
}

// NewVariablePool creates a new variable pool
func NewVariablePool(ctx *GeneratorContext) *VariablePool {
	return &VariablePool{
		ctx:   ctx,
		cache: make(map[string]*SharedVariables),
	}
}

// GetSharedVariables gets shared variables for a specific category (Flyweight Object)
// Returns cached instance if exists, otherwise creates a new one
func (p *VariablePool) GetSharedVariables(category string) *SharedVariables {
	p.mu.RLock()
	if cached, exists := p.cache[category]; exists {
		p.mu.RUnlock()
		return cached
	}
	p.mu.RUnlock()

	// Create new shared variable set
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check
	if cached, exists := p.cache[category]; exists {
		return cached
	}

	shared := p.createSharedVariables(category)
	shared.Freeze() // Freeze to prevent modification
	p.cache[category] = shared
	return shared
}

// createSharedVariables creates shared variables based on category
func (p *VariablePool) createSharedVariables(category string) *SharedVariables {
	shared := &SharedVariables{
		category: category,
		vars:     make(map[string]interface{}),
	}

	switch category {
	case CategoryCommon:
		p.fillCommonVariables(shared)
	case CategoryBuild:
		p.fillBuildVariables(shared)
	case CategoryRuntime:
		p.fillRuntimeVariables(shared)
	case CategoryPlugin:
		p.fillPluginVariables(shared)
	case CategoryCIPaths:
		p.fillCIPathVariables(shared)
	case CategoryService:
		p.fillServiceVariables(shared)
	case CategoryLanguage:
		p.fillLanguageVariables(shared)
	}

	return shared
}

// fillCommonVariables fills common variables
func (p *VariablePool) fillCommonVariables(shared *SharedVariables) {
	cfg := p.ctx.Config
	shared.vars[VarServiceName] = cfg.Service.Name
	shared.vars[VarDeployDir] = cfg.Service.DeployDir
	shared.vars[VarServiceRoot] = p.ctx.Paths.ServiceRoot
	shared.vars[VarGeneratedAt] = cfg.Metadata.GeneratedAt
	shared.vars[VarConfigDir] = p.ctx.Paths.ConfigDir
	shared.vars[VarServiceBinDir] = p.ctx.Paths.BinDir
}

// fillBuildVariables fills build-related variables
func (p *VariablePool) fillBuildVariables(shared *SharedVariables) {
	cfg := p.ctx.Config
	shared.vars[VarBuildCommand] = cfg.Build.Commands.Build
	shared.vars[VarPreBuildCommand] = cfg.Build.Commands.PreBuild
	shared.vars[VarPostBuildCommand] = cfg.Build.Commands.PostBuild
	shared.vars["BUILD_DEPS_PACKAGES"] = cfg.Build.SystemDependencies.Packages
	shared.vars["BUILDER_IMAGE_AMD64"] = cfg.Build.BuilderImage.AMD64
	shared.vars["BUILDER_IMAGE_ARM64"] = cfg.Build.BuilderImage.ARM64
	shared.vars["RUNTIME_IMAGE_AMD64"] = cfg.Build.RuntimeImage.AMD64
	shared.vars["RUNTIME_IMAGE_ARM64"] = cfg.Build.RuntimeImage.ARM64
}

// fillRuntimeVariables fills runtime-related variables
func (p *VariablePool) fillRuntimeVariables(shared *SharedVariables) {
	cfg := p.ctx.Config
	shared.vars["STARTUP_COMMAND"] = cfg.Runtime.Startup.Command
	shared.vars["ENV_VARS"] = cfg.Runtime.Startup.Env
	shared.vars["RUNTIME_DEPS_PACKAGES"] = cfg.Runtime.SystemDependencies.Packages
	shared.vars["HEALTHCHECK_ENABLED"] = cfg.Runtime.Healthcheck.Enabled
	shared.vars["HEALTHCHECK_TYPE"] = cfg.Runtime.Healthcheck.Type
	shared.vars["GENERATE_SCRIPTS"] = cfg.Runtime.GenerateScripts
}

// fillPluginVariables fills plugin-related variables
func (p *VariablePool) fillPluginVariables(shared *SharedVariables) {
	cfg := p.ctx.Config
	shared.vars[VarPluginRootDir] = DefaultPluginRootDir
	shared.vars["PLUGIN_INSTALL_DIR"] = cfg.Plugins.InstallDir
	shared.vars["HAS_PLUGINS"] = len(cfg.Plugins.Items) > 0

	// Pre-process plugin information
	var plugins []map[string]interface{}
	for _, plugin := range cfg.Plugins.Items {
		plugins = append(plugins, map[string]interface{}{
			"Name":        plugin.Name,
			"DownloadURL": plugin.DownloadURL,
			"InstallDir":  cfg.Plugins.InstallDir,
			"RuntimeEnv":  plugin.RuntimeEnv,
		})
	}
	shared.vars["PLUGINS"] = plugins
}

// fillCIPathVariables fills CI path variables
func (p *VariablePool) fillCIPathVariables(shared *SharedVariables) {
	if p.ctx.Paths.CI != nil {
		for k, v := range p.ctx.Paths.CI.ToTemplateVars() {
			shared.vars[k] = v
		}
	}
}

// fillServiceVariables fills service-related variables
func (p *VariablePool) fillServiceVariables(shared *SharedVariables) {
	cfg := p.ctx.Config
	shared.vars["PORTS"] = cfg.Service.Ports

	if len(cfg.Service.Ports) > 0 {
		shared.vars[VarServicePort] = cfg.Service.Ports[0].Port
	}

	// Expose ports
	var exposePorts []int
	for _, port := range cfg.Service.Ports {
		if port.Expose {
			exposePorts = append(exposePorts, port.Port)
		}
	}
	shared.vars["EXPOSE_PORTS"] = exposePorts
}

// fillLanguageVariables fills language-related variables
func (p *VariablePool) fillLanguageVariables(shared *SharedVariables) {
	cfg := p.ctx.Config
	shared.vars[VarLanguage] = cfg.Language.Type
	shared.vars[VarLanguageVersion] = cfg.Language.Version
	shared.vars["LANGUAGE_CONFIG"] = cfg.Language.Config
}

// Freeze freezes the variable set to prevent modification
func (s *SharedVariables) Freeze() {
	s.frozen = true
}

// IsFrozen returns whether the variable set is frozen
func (s *SharedVariables) IsFrozen() bool {
	return s.frozen
}

// ToMap returns a copy of variables (prevents external modification)
func (s *SharedVariables) ToMap() map[string]interface{} {
	result := make(map[string]interface{}, len(s.vars))
	for k, v := range s.vars {
		result[k] = v
	}
	return result
}

// Get gets a single variable
func (s *SharedVariables) Get(key string) (interface{}, bool) {
	val, exists := s.vars[key]
	return val, exists
}

// Category returns the category of this shared variable set
func (s *SharedVariables) Category() string {
	return s.category
}

// Size returns the number of variables in this set
func (s *SharedVariables) Size() int {
	return len(s.vars)
}
