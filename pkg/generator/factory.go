package generator

import (
	"fmt"
	"sync"

	"github.com/junjiewwang/service-template/pkg/config"
)

// TemplateGenerator defines the interface for all template generators
type TemplateGenerator interface {
	// Generate generates content from template
	Generate() (string, error)
	// GetName returns the generator name
	GetName() string
}

// GeneratorCreator is a function type that creates a template generator
type GeneratorCreator func(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error)

// generatorRegistry holds all registered generator creators
var (
	generatorRegistry = make(map[string]GeneratorCreator)
	registryMutex     sync.RWMutex
)

// RegisterGenerator registers a generator creator with a unique type identifier
func RegisterGenerator(generatorType string, creator GeneratorCreator) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	if _, exists := generatorRegistry[generatorType]; exists {
		panic(fmt.Sprintf("generator type %s is already registered", generatorType))
	}
	generatorRegistry[generatorType] = creator
}

// GetRegisteredGenerators returns all registered generator types
func GetRegisteredGenerators() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	types := make([]string, 0, len(generatorRegistry))
	for t := range generatorRegistry {
		types = append(types, t)
	}
	return types
}

// GeneratorFactory creates template generators
type GeneratorFactory struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
}

// NewGeneratorFactory creates a new generator factory
func NewGeneratorFactory(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *GeneratorFactory {
	return &GeneratorFactory{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
	}
}

// CreateGenerator creates a generator by type using the registered creator
func (f *GeneratorFactory) CreateGenerator(generatorType string, options ...interface{}) (TemplateGenerator, error) {
	registryMutex.RLock()
	creator, exists := generatorRegistry[generatorType]
	registryMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown generator type: %s (available types: %v)", generatorType, GetRegisteredGenerators())
	}

	return creator(f.config, f.templateEngine, f.variables, options...)
}

// BaseTemplateGenerator provides common functionality for all generators
type BaseTemplateGenerator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
	name           string
}

// GetName returns the generator name
func (g *BaseTemplateGenerator) GetName() string {
	return g.name
}

// RenderTemplate renders a template with variables
func (g *BaseTemplateGenerator) RenderTemplate(template string, vars map[string]interface{}) (string, error) {
	return g.templateEngine.Render(template, vars)
}
