package core

import (
	"github.com/junjiewwang/service-template/pkg/generator/context"
)

// BaseGenerator provides common functionality for all generators
type BaseGenerator struct {
	name   string
	ctx    *context.GeneratorContext
	engine *TemplateEngine
}

// NewBaseGenerator creates a new base generator
func NewBaseGenerator(name string, ctx *context.GeneratorContext, engine *TemplateEngine) BaseGenerator {
	return BaseGenerator{
		name:   name,
		ctx:    ctx,
		engine: engine,
	}
}

// GetName returns the generator name
func (g *BaseGenerator) GetName() string {
	return g.name
}

// GetContext returns the generator context
func (g *BaseGenerator) GetContext() *context.GeneratorContext {
	return g.ctx
}

// Validate validates the generator configuration (default implementation)
func (g *BaseGenerator) Validate() error {
	if g.ctx == nil {
		return ErrNilContext
	}
	if g.ctx.Config == nil {
		return ErrNilConfig
	}
	return nil
}

// RenderTemplate renders a template with variables
func (g *BaseGenerator) RenderTemplate(template string, vars map[string]interface{}) (string, error) {
	return g.engine.Render(template, vars)
}

// RenderTemplateWithName renders a named template with variables
func (g *BaseGenerator) RenderTemplateWithName(name, template string, vars map[string]interface{}) (string, error) {
	return g.engine.RenderWithName(name, template, vars)
}
