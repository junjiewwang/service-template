package context

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// GeneratorContext holds all context information needed for generation
type GeneratorContext struct {
	// Config is the service configuration
	Config *config.ServiceConfig

	// Paths manages all path-related information
	Paths *Paths

	// OutputDir is the output directory for generated files
	OutputDir string

	// VariablePool manages shared variables (Flyweight Pattern)
	VariablePool *VariablePool
}

// NewGeneratorContext creates a new generator context
func NewGeneratorContext(cfg *config.ServiceConfig, outputDir string) *GeneratorContext {
	paths := NewPaths(cfg)

	ctx := &GeneratorContext{
		Config:    cfg,
		Paths:     paths,
		OutputDir: outputDir,
	}

	// Initialize variable pool
	ctx.VariablePool = NewVariablePool(ctx)

	return ctx
}

// GetVariableComposer returns a new variable composer
func (c *GeneratorContext) GetVariableComposer() *VariableComposer {
	return NewVariableComposer(c.VariablePool)
}

// GetVariablePreset returns a new variable preset
func (c *GeneratorContext) GetVariablePreset() *VariablePreset {
	return NewVariablePreset(c.VariablePool)
}

// Validate validates the context
func (c *GeneratorContext) Validate() error {
	if c.Config == nil {
		return ErrNilConfig
	}
	if c.Paths == nil {
		return ErrNilPaths
	}
	if c.VariablePool == nil {
		return ErrNilVariablePool
	}
	return nil
}
