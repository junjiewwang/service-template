package core

import (
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/domain/events"
)

// BaseGenerator provides common functionality for all generators
type BaseGenerator struct {
	name      string
	ctx       *context.GeneratorContext
	engine    *TemplateEngine
	publisher events.EventPublisher
}

// NewBaseGenerator creates a new base generator
func NewBaseGenerator(name string, ctx *context.GeneratorContext, engine *TemplateEngine) BaseGenerator {
	return BaseGenerator{
		name:      name,
		ctx:       ctx,
		engine:    engine,
		publisher: events.NewSimpleEventPublisher(), // Default publisher
	}
}

// NewBaseGeneratorWithPublisher creates a new base generator with a custom event publisher
func NewBaseGeneratorWithPublisher(name string, ctx *context.GeneratorContext, engine *TemplateEngine, publisher events.EventPublisher) BaseGenerator {
	return BaseGenerator{
		name:      name,
		ctx:       ctx,
		engine:    engine,
		publisher: publisher,
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

// GetEngine returns the template engine
func (g *BaseGenerator) GetEngine() *TemplateEngine {
	return g.engine
}

// GetPublisher returns the event publisher
func (g *BaseGenerator) GetPublisher() events.EventPublisher {
	return g.publisher
}

// SetPublisher sets the event publisher
func (g *BaseGenerator) SetPublisher(publisher events.EventPublisher) {
	g.publisher = publisher
}

// PublishEvent publishes an event
func (g *BaseGenerator) PublishEvent(event events.Event) error {
	if g.publisher == nil {
		return nil // No publisher, silently ignore
	}
	return g.publisher.Publish(event)
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

// VariablePreparator is an interface for generators that need custom variable preparation
type VariablePreparator interface {
	// PrepareCustomVariables prepares generator-specific custom variables
	// It receives a composer and can add/override variables as needed
	PrepareCustomVariables(composer *context.VariableComposer) error
}

// PrepareVariables provides a standard way to prepare variables for templates
// It uses the variable preset system and allows generators to add custom variables
func (g *BaseGenerator) PrepareVariables(presetFunc func() *context.VariableComposer) (map[string]interface{}, error) {
	// Get the preset composer
	composer := presetFunc()

	// If the generator implements VariablePreparator, call it
	if preparator, ok := interface{}(g).(VariablePreparator); ok {
		if err := preparator.PrepareCustomVariables(composer); err != nil {
			return nil, err
		}
	}

	return composer.Build(), nil
}
