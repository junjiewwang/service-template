package context

// VariableComposer composes different categories of shared variables (Flyweight Pattern Client)
// Provides a fluent API for building variable sets
type VariableComposer struct {
	pool   *VariablePool
	result map[string]interface{}
}

// NewVariableComposer creates a new variable composer
func NewVariableComposer(pool *VariablePool) *VariableComposer {
	return &VariableComposer{
		pool:   pool,
		result: make(map[string]interface{}),
	}
}

// WithCommon adds common variables
func (c *VariableComposer) WithCommon() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryCommon))
	return c
}

// WithBuild adds build variables
func (c *VariableComposer) WithBuild() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryBuild))
	return c
}

// WithRuntime adds runtime variables
func (c *VariableComposer) WithRuntime() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryRuntime))
	return c
}

// WithPlugin adds plugin variables
func (c *VariableComposer) WithPlugin() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryPlugin))
	return c
}

// WithCIPaths adds CI path variables
func (c *VariableComposer) WithCIPaths() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryCIPaths))
	return c
}

// WithService adds service variables
func (c *VariableComposer) WithService() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryService))
	return c
}

// WithLanguage adds language variables
func (c *VariableComposer) WithLanguage() *VariableComposer {
	c.merge(c.pool.GetSharedVariables(CategoryLanguage))
	return c
}

// WithAll adds all standard variable categories
func (c *VariableComposer) WithAll() *VariableComposer {
	return c.
		WithCommon().
		WithBuild().
		WithRuntime().
		WithPlugin().
		WithCIPaths().
		WithService().
		WithLanguage()
}

// WithArchitecture adds architecture-specific variables (extrinsic state)
func (c *VariableComposer) WithArchitecture(arch string) *VariableComposer {
	c.result[VarGOARCH] = arch
	c.result["ARCH"] = arch

	if arch == "amd64" || arch == "arm64" {
		c.result[VarGOOS] = "linux"
	}

	// Select architecture-specific images from shared variables
	buildVars := c.pool.GetSharedVariables(CategoryBuild)
	switch arch {
	case "amd64":
		if img, ok := buildVars.Get("BUILDER_IMAGE_AMD64"); ok {
			c.result["BUILDER_IMAGE"] = img
		}
		if img, ok := buildVars.Get("RUNTIME_IMAGE_AMD64"); ok {
			c.result["RUNTIME_IMAGE"] = img
		}
	case "arm64":
		if img, ok := buildVars.Get("BUILDER_IMAGE_ARM64"); ok {
			c.result["BUILDER_IMAGE"] = img
		}
		if img, ok := buildVars.Get("RUNTIME_IMAGE_ARM64"); ok {
			c.result["RUNTIME_IMAGE"] = img
		}
	}

	return c
}

// WithCustom adds a custom variable (extrinsic state)
func (c *VariableComposer) WithCustom(key string, value interface{}) *VariableComposer {
	c.result[key] = value
	return c
}

// WithCustomMap adds multiple custom variables in batch
func (c *VariableComposer) WithCustomMap(vars map[string]interface{}) *VariableComposer {
	for k, v := range vars {
		c.result[k] = v
	}
	return c
}

// Override overrides an existing variable
func (c *VariableComposer) Override(key string, value interface{}) *VariableComposer {
	c.result[key] = value
	return c
}

// OverrideMap overrides multiple variables in batch
func (c *VariableComposer) OverrideMap(vars map[string]interface{}) *VariableComposer {
	for k, v := range vars {
		c.result[k] = v
	}
	return c
}

// Has checks if a variable exists
func (c *VariableComposer) Has(key string) bool {
	_, exists := c.result[key]
	return exists
}

// Get gets a variable value
func (c *VariableComposer) Get(key string) (interface{}, bool) {
	val, exists := c.result[key]
	return val, exists
}

// Build builds the final variable set
func (c *VariableComposer) Build() map[string]interface{} {
	return c.result
}

// Clone creates a copy of the composer
func (c *VariableComposer) Clone() *VariableComposer {
	newComposer := &VariableComposer{
		pool:   c.pool,
		result: make(map[string]interface{}, len(c.result)),
	}
	for k, v := range c.result {
		newComposer.result[k] = v
	}
	return newComposer
}

// Size returns the number of variables
func (c *VariableComposer) Size() int {
	return len(c.result)
}

// merge merges shared variables into the result set
func (c *VariableComposer) merge(shared *SharedVariables) {
	for k, v := range shared.ToMap() {
		// Don't override existing variables
		if _, exists := c.result[k]; !exists {
			c.result[k] = v
		}
	}
}

// VariablePreset provides preset variable combinations for common scenarios
type VariablePreset struct {
	composer *VariableComposer
}

// NewVariablePreset creates a new variable preset
func NewVariablePreset(pool *VariablePool) *VariablePreset {
	return &VariablePreset{
		composer: NewVariableComposer(pool),
	}
}

// ForDockerfile returns a preset for Dockerfile generation
func (p *VariablePreset) ForDockerfile(arch string) *VariableComposer {
	return p.composer.Clone().
		WithCommon().
		WithBuild().
		WithRuntime().
		WithPlugin().
		WithCIPaths().
		WithService().
		WithLanguage().
		WithArchitecture(arch)
}

// ForBuildScript returns a preset for build script generation
func (p *VariablePreset) ForBuildScript() *VariableComposer {
	return p.composer.Clone().
		WithCommon().
		WithBuild().
		WithPlugin().
		WithCIPaths()
}

// ForCompose returns a preset for docker-compose generation
func (p *VariablePreset) ForCompose() *VariableComposer {
	return p.composer.Clone().
		WithCommon().
		WithRuntime().
		WithService()
}

// ForMakefile returns a preset for Makefile generation
func (p *VariablePreset) ForMakefile() *VariableComposer {
	return p.composer.Clone().
		WithCommon().
		WithService().
		WithCIPaths()
}

// ForDevOps returns a preset for DevOps configuration generation
func (p *VariablePreset) ForDevOps() *VariableComposer {
	return p.composer.Clone().
		WithCommon().
		WithBuild().
		WithLanguage()
}

// ForScript returns a preset for script generation
func (p *VariablePreset) ForScript() *VariableComposer {
	return p.composer.Clone().
		WithCommon().
		WithLanguage().
		WithRuntime().
		WithService().
		WithCIPaths()
}
